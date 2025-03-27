package main

import (
	"encoding/json" // Для работы с JSON
	"fmt"           // Для вывода на экран
	"io/ioutil"     // Для работы с файлами
	"log"           // Для логирования ошибок
	"os"            // Для работы с операционной системой
)

const fileName = "tasks.json" // Имя файла для хранения задач

// Структура задачи
type Task struct {
	ID       int    `json:"id"`       // Уникальный идентификатор задачи
	Title    string `json:"title"`    // Название задачи
	Complete bool   `json:"complete"` // Статус выполнения задачи
}

// Функция для чтения задач из файла
func readTasksFromFile() ([]Task, error) {
	var tasks []Task

	// Чтение содержимого файла
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// Если файл не существует, возвращаем пустой список задач
			return tasks, nil
		}
		return nil, err // Возвращаем ошибку, если не удалось прочитать файл
	}

	// Декодируем JSON в структуру
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err // Возвращаем ошибку, если не удалось распарсить JSON
	}

	return tasks, nil
}

// Функция для записи задач в файл
func writeTasksToFile(tasks []Task) error {
	// Преобразуем задачи в JSON
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err // Возвращаем ошибку, если не удалось преобразовать в JSON
	}

	// Записываем данные в файл
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err // Возвращаем ошибку, если не удалось записать в файл
	}

	return nil
}

// Функция для добавления задачи
func addTask(title string) {
	tasks, err := readTasksFromFile()
	if err != nil {
		log.Fatal("Ошибка при чтении задач:", err)
	}

	// Создаём новую задачу
	id := len(tasks) + 1 // Новый ID для задачи (просто инкрементируем количество задач)
	task := Task{ID: id, Title: title, Complete: false}

	// Добавляем задачу в список
	tasks = append(tasks, task)

	// Записываем обновлённый список в файл
	err = writeTasksToFile(tasks)
	if err != nil {
		log.Fatal("Ошибка при записи задач:", err)
	}

	fmt.Println("Задача добавлена:", task)
}

// Функция для отображения всех задач
func listTasks() {
	tasks, err := readTasksFromFile()
	if err != nil {
		log.Fatal("Ошибка при чтении задач:", err)
	}

	if len(tasks) == 0 {
		fmt.Println("Нет задач для отображения.")
		return
	}

	fmt.Println("Список задач:")
	for _, task := range tasks {
		status := "Не выполнено"
		if task.Complete {
			status = "Выполнено"
		}
		fmt.Printf("%d. %s [%s]\n", task.ID, task.Title, status)
	}
}

// Функция для отметки задачи как выполненной
func markTaskComplete(id int) {
	tasks, err := readTasksFromFile()
	if err != nil {
		log.Fatal("Ошибка при чтении задач:", err)
	}

	// Ищем задачу по ID
	var taskIndex = -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		fmt.Println("Задача с таким ID не найдена.")
		return
	}

	// Отмечаем задачу как выполненную
	tasks[taskIndex].Complete = true

	// Записываем обновлённый список в файл
	err = writeTasksToFile(tasks)
	if err != nil {
		log.Fatal("Ошибка при записи задач:", err)
	}

	fmt.Println("Задача отмечена как выполненная:", tasks[taskIndex])
}

func main() {
	// Главное меню
	for {
		fmt.Println("\nМенеджер задач - To-Do List")
		fmt.Println("1. Показать все задачи")
		fmt.Println("2. Добавить задачу")
		fmt.Println("3. Отметить задачу как выполненную")
		fmt.Println("4. Выйти")

		var choice int
		fmt.Print("Выберите опцию: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			log.Fatal("Ошибка ввода:", err)
		}

		switch choice {
		case 1:
			listTasks()
		case 2:
			var title string
			fmt.Print("Введите название задачи: ")
			fmt.Scanln(&title)
			addTask(title)
		case 3:
			var id int
			fmt.Print("Введите ID задачи для пометки как выполненную: ")
			fmt.Scanln(&id)
			markTaskComplete(id)
		case 4:
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Неверный выбор. Пожалуйста, выберите опцию снова.")
		}
	}
}
