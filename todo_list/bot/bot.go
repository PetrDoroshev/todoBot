package bot

import (
	"fmt"
	"github.com/PetrDoroshev/HomeWork_db/golang/todo_list/data_base"
	_ "github.com/Syfaro/telegram-bot-api"
	"github.com/go-co-op/gocron"
	_ "github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
	"strings"
	"time"
)

const layout = "2006-01-02 15:04"

type State int

const (
	EnterDescription  State = 1
	EnterPriority           = 2
	EnterDate               = 3
	EnterNumberOfTask       = 4
)

type Bot struct {
	todoList      *data_base.DB
	api           *tgbotapi.BotAPI
	userState     map[int64]State
	userTask      map[int64]*data_base.Task
	notifiedTasks []data_base.Task
}

func NewBot(TodoList *data_base.DB) *Bot {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))

	if err != nil {
		panic(err)
	}
	return &Bot{todoList: TodoList,
		api:           bot,
		userState:     make(map[int64]State),
		userTask:      make(map[int64]*data_base.Task),
		notifiedTasks: make([]data_base.Task, 0, 0)}
}

func (bot *Bot) BotLoop() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.api.GetUpdatesChan(u)

	s := gocron.NewScheduler(time.UTC)
	_, _ = s.Every(1).Minute().Do(bot.checkTasks)
	s.StartAsync()

	if err != nil {
		panic(err)
	}

	for update := range updates {

		if update.Message != nil {

			m := update.Message.Text
			if m == "/add_task" {
				bot.userTask[update.Message.Chat.ID] = &data_base.Task{UserID: update.Message.Chat.ID}
				bot.userState[update.Message.Chat.ID] = EnterDescription
				bot.sendMessage(update.Message.Chat.ID, "Enter description of task: \n"+
					"Enter /cancel command to cancel creating task")

			} else if m == "/finish_task" {
				bot.finishTask(update.Message.Chat.ID)

			} else if update.Message.Command() == "show_tasks" {
				bot.showTasks(update.Message.Chat.ID, m)

			} else if m == "/show_finishedtasks" {
				bot.showFinishedTasks(update.Message.Chat.ID)
			} else if m == "/show_stats" {
				ft, _ := bot.todoList.CountFinishedTasks(update.Message.Chat.ID)
				bot.sendMessage(update.Message.Chat.ID, fmt.Sprintf("Statistics: \n\nfinished_tasks: %d", ft))

			} else if m == "/cancel" {
				delete(bot.userTask, update.Message.Chat.ID)
				delete(bot.userState, update.Message.Chat.ID)

			} else if bot.userState[update.Message.Chat.ID] == EnterDescription {
				bot.enterDescriptionDialog(update.Message.Chat.ID, m)

			} else if bot.userState[update.Message.Chat.ID] == EnterPriority {
				bot.enterPriorityDialog(update.Message.Chat.ID, m)

			} else if bot.userState[update.Message.Chat.ID] == EnterDate {
				bot.enterDateDialog(update.Message.Chat.ID, m)

			} else if bot.userState[update.Message.Chat.ID] == EnterNumberOfTask {
				n, err := strconv.Atoi(m)
				if err == nil {
					err = bot.todoList.FinishTask(update.Message.Chat.ID, n)
					if err != nil {
						bot.sendMessage(update.Message.Chat.ID, "Error occurred while updating task status")
					} else {
						delete(bot.userState, update.Message.Chat.ID)
						bot.sendMessage(update.Message.Chat.ID, "Task finished")
					}
				} else {
					bot.sendMessage(update.Message.Chat.ID, "Not a number")
				}
			}
		} else if update.CallbackQuery != nil {

			_, ok := bot.userTask[update.CallbackQuery.Message.Chat.ID]

			if update.CallbackQuery.Data == "in work" && ok {

				bot.userTask[update.CallbackQuery.Message.Chat.ID].Status = "in work"
				err := bot.todoList.AddNewTask(bot.userTask[update.CallbackQuery.Message.Chat.ID])

				if bot.userTask[update.CallbackQuery.Message.Chat.ID].NotifyTime != nil {
					bot.notifiedTasks = append(bot.notifiedTasks, *bot.userTask[update.CallbackQuery.Message.Chat.ID])
				}

				if err != nil {
					bot.sendMessage(update.Message.Chat.ID, "Error occurred while adding task to list")
				} else {
					bot.sendMessage(update.CallbackQuery.Message.Chat.ID, "Task successfully added to list")
					delete(bot.userTask, update.CallbackQuery.Message.Chat.ID)
				}

			}
		}
	}
}

func (bot *Bot) checkTasks() {

	now := time.Now().Format(layout)
	for i, task := range bot.notifiedTasks {
		if *task.NotifyTime == now {
			bot.sendMessage(task.UserID, fmt.Sprintf("You have unfisnished task: %s", task.Description))
			bot.notifiedTasks[i] = bot.notifiedTasks[len(bot.notifiedTasks)-1]
			bot.notifiedTasks = bot.notifiedTasks[:len(bot.notifiedTasks)-1]
		}
	}
}

func (bot *Bot) showFinishedTasks(chatID int64) {
	tasks, _ := bot.todoList.ListFinishedTasks(chatID)
	msg := "Finished tasks: \n\n"
	for i, task := range tasks {
		msg += fmt.Sprintf("%d. %s\n", i+1, task.Description)
	}
	bot.sendMessage(chatID, msg)
}

func (bot *Bot) showTasks(chatID int64, m string) {
	var tasks []data_base.Task
	s := strings.Split(m+" ", " ")

	if s[1] == "low" || s[1] == "medium" || s[1] == "high" {
		tasks, _ = bot.todoList.ListTaskWithPriority(chatID, s[1])
	} else {
		tasks, _ = bot.todoList.ListUserTasks(chatID)
	}
	msg := "Todo list: \n\n"
	for i, task := range tasks {
		msg += fmt.Sprintf("%d. %s\n\t\t\t priority: %s\n", i+1, task.Description, task.Priority)
		if task.NotifyTime != nil {
			t, _ := time.Parse("2006-01-02T15:04:05Z", *task.NotifyTime)
			msg += "\t\t\t notify date: " + t.Format(layout) + "\n"
		}
		msg += "\n"
	}
	bot.sendMessage(chatID, msg)

}

func (bot *Bot) finishTask(chatID int64) {
	tasks, _ := bot.todoList.ListUserTasks(chatID)
	msg := "Todo list: \n"
	for i, task := range tasks {
		msg += fmt.Sprintf("%d. %s\n", i+1, task.Description)
	}
	msg += "Enter number of task to finish: "
	bot.sendMessage(chatID, msg)
	bot.userState[chatID] = EnterNumberOfTask
}

func (bot *Bot) enterDescriptionDialog(chatID int64, m string) {
	bot.userTask[chatID].Description = m
	bot.userState[chatID] = EnterPriority
	bot.sendMessage(chatID, "Now, enter priority of the task (low, medium or high)")
}

func (bot *Bot) enterPriorityDialog(chatID int64, m string) {
	if m == "low" || m == "medium" || m == "high" {
		bot.userTask[chatID].Priority = m
		bot.userState[chatID] = EnterDate
		bot.sendMessage(chatID, "Enter date of notification in YYYY-MM-DD HH:MM format or \"skip\": ")
	} else {
		bot.sendMessage(chatID, "Incorrect value of priority")
	}
}

func (bot *Bot) enterDateDialog(chatID int64, m string) {
	if m == "skip" || checkDateFormat(m) {
		if m != "skip" {
			bot.userTask[chatID].NotifyTime = &m
		}
		delete(bot.userState, chatID)

		nt := ""
		if bot.userTask[chatID].NotifyTime != nil {
			nt = *bot.userTask[chatID].NotifyTime
		}

		text := fmt.Sprintf("%s\nPritority: %s\nNotify time: %s\nStatus: %s\n",
			bot.userTask[chatID].Description,
			bot.userTask[chatID].Priority,
			nt,
			bot.userTask[chatID].Status)

		button := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Take this task", "in work"),
		))
		bot.sendMessageWithMarkup(chatID, text, button)

	} else {
		bot.sendMessage(chatID, "incorrect date format")
	}
}

func (bot *Bot) sendMessage(chatID int64, message string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.api.Send(msg)

	return msg
}

func (bot *Bot) sendMessageWithMarkup(chatID int64, message string, buttons tgbotapi.InlineKeyboardMarkup) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = buttons
	bot.api.Send(msg)

	return msg
}

func checkDateFormat(str string) bool {
	const layout = "2006-01-02 15:04"
	_, err := time.Parse(layout, str)

	return err == nil
}
