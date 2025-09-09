package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// ----------------------------------Правила сортировки------------------------------------
type GetRules interface {
	SortFiles(workingDir string, listOfFiles []os.DirEntry) map[string]string
	Reminder(workingDir string, logFile *os.File, pkgs map[string]string)
}
type SortByExtension struct {
}
type SortByContent struct {
}
type actions struct {
	file        string
	source      string
	destination string
}

// --------------------------Функция для отслеживания переименования папок по умолчанию-----------------
func renameDefPkg(workingDir string) {}

// ---------------------------------Функция для напоминаний-------------------------------
func MakeMap(pkgs map[string]string) map[string][]string {
	revPkgs := make(map[string][]string)
	for extension, packageName := range pkgs {
		revPkgs[packageName] = append(revPkgs[packageName], extension)
	}
	return revPkgs
}
func ApplyRules() {
	/**/
}
func getTrgetPckg(workingDir string, pkgs map[string]string, revPkgs map[string][]string) {
	elems, _ := os.ReadDir(workingDir)
	for _, elem := range elems {
		if elem.IsDir() {
			content, _ := os.ReadDir(path.Join(workingDir, elem.Name()))
			targetPackage := ""
			for _, files := range content {
				extension := strings.ToLower(path.Ext(files.Name()))
				if packageName, ok := pkgs[extension]; ok {
					targetPackage = packageName
					break
				}
			}
			if targetPackage != "" {
				allExts := revPkgs[targetPackage]
				fmt.Printf("Папка %s содержит файлы категории %s (%v)\n", elem.Name(), targetPackage, allExts)
			}
		}
	}
}
func (ct SortByContent) Reminder(workingDir string, logFile *os.File, pkgs map[string]string) {
	revPkgs := MakeMap(pkgs)
	getTrgetPckg(workingDir, pkgs, revPkgs)
}

/*
Для изображений - более 2х недель или более 5Мб
Для видео - более 2х недель (или флаг просмотра в будущем - реализация собственного плеера)
Для книг - более года (или флаг просмотра в будущем - реализация своего ридера)
Для таблиц - более месяца
Для музыки - более 2х недель или особый флаг
Для текстовых файлов doc и docx - более месяца, остальное - более года
Для архивов - более полугода
Презентации и пдф - более 2х недель
*/
// --------------------------------Функции для логирования-------------------------------
func LogPack(pkgs map[string]string) (string, error) {
	nameLogPack, ok := pkgs[".log"]
	if ok {
		return nameLogPack, nil
	} else {
		nameLogPack = ""
		errText := "Предупреждение! Не назначена папка для сохранения логов! Логи сохраняются в папку по умолчанию.\n"
		return nameLogPack, errors.New(errText)
	}
}
func MakeLogFile(err error, workingDir string, nameLogPack string) *os.File {
	fmt.Println(err)
	if err == nil {
		logFile, errCrLogF := os.OpenFile(path.Join(workingDir, nameLogPack, "log-"+time.Now().Format("2006-01-02")+".log"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		if errCrLogF != nil {
			fmt.Println(errCrLogF.Error())
		}
		return logFile
	} else {
		os.MkdirAll(workingDir+nameLogPack, 0755)
		logFile, errCrLogF := os.OpenFile(path.Join(workingDir, nameLogPack, "log-"+time.Now().Format("2006-01-02")+".log"), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if errCrLogF != nil {
			fmt.Println(errCrLogF.Error())
		}
		logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": " + err.Error())
		return logFile
	}
}

func WriteLog(file string, packageName string, workingDir string, logFile *os.File) {
	textLog := time.Now().Format("2006-01-02 15:04:05") + ": File " + file + " was removed to " + packageName
	logFile.WriteString(textLog + "\n")
}

// -------------------------------Функция для отмены предыдущего действия--------------------------------
func Undo(undoStack *[]actions, redoStack *[]actions, logFile *os.File) error {
	if len(*undoStack) != 0 {
		action := (*undoStack)[len(*undoStack)-1]
		*undoStack = (*undoStack)[:len(*undoStack)-1]
		currPath := path.Join(action.source, action.destination, action.file)
		revPath := path.Join(action.source, action.file)
		os.Rename(currPath, revPath)
		logFile.WriteString("----UNDO:----")
		WriteLog(action.file, action.source, action.source, logFile)
		*redoStack = append(*redoStack, action)
		return nil
	} else {
		errText := "Ошибка! Нет действий для отмены!\n"
		logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": " + errText)
		return errors.New(errText)
	}
}

// -------------------------------Функция для возврата предыдущего действия------------------------------
func Redo(undoStack *[]actions, redoStack *[]actions, logFile *os.File) error {
	if len(*redoStack) != 0 {
		action := (*redoStack)[len(*redoStack)-1]
		*redoStack = (*redoStack)[:len(*redoStack)-1]
		currPath := path.Join(action.source, action.file)
		newPath := path.Join(action.source, action.destination, action.file)
		os.Rename(currPath, newPath)
		logFile.WriteString("----REDO:----")
		WriteLog(action.file, action.destination, action.source, logFile)
		*undoStack = append(*undoStack, action)
		return nil
	} else {
		errText := "Ошибка! Нет отмененных действий для возврата!\n"
		logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": " + errText)
		return errors.New(errText)
	}
}

//--------------------------------Вспомогательные функции для сортировки--------------------------------

func Sort(pkgs map[string]string, workingDir string, listOfFiles []os.DirEntry, logFile *os.File) ([]actions, []actions) { //Разделить
	//Создание папок и стэков - MkFilesNStakes
	//SortAlg - непосредственно сортировка
	//
	undoStack := make([]actions, 0)
	redoStack := make([]actions, 0)
	for _, packageName := range pkgs {
		os.Mkdir(workingDir+packageName, 0755)
	}
	for _, file := range listOfFiles {
		extensions := strings.ToLower(path.Ext(file.Name()))
		if !(file.IsDir()) && extensions != ".log" {
			packageName, ok := pkgs[extensions]
			if !ok {
				os.Mkdir(workingDir+"Others", 0755)
				packageName = "Others"
			}
			file := file.Name()
			os.Rename(path.Join(workingDir+file), path.Join(workingDir, packageName, file))
			WriteLog(file, packageName, workingDir, logFile)
			action := actions{file, workingDir, packageName}
			undoStack = append(undoStack, action)
		}
	}
	return undoStack, redoStack
}
func DelOldLogs(workingDir string, logFile *os.File, nameLogPack string) {
	allLogs, errScanDir := os.ReadDir(path.Join(workingDir, nameLogPack))
	if errScanDir != nil {
		logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": Ошибка сканирования директории!\n")
	} else {
		for _, log := range allLogs {
			if path.Ext(log.Name()) == ".log" {
				logStat, errGetInf := os.Stat(path.Join(workingDir, nameLogPack, log.Name()))
				if errGetInf != nil {
					logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": Ошибка проверки информации о логах!\n")
					logFile.WriteString("Функция удаления устаревших логов отключена!\n")
				} else {
					formatedLogInfo := logStat.ModTime()
					fmt.Println(formatedLogInfo.Format("2006-01-02 15:04:05"))
					if time.Since(formatedLogInfo) > 72*time.Hour {
						os.Remove(path.Join(workingDir, nameLogPack, log.Name()))
						logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": Устаревшие логи были удалены: " + log.Name() + "\n")

					}

				}
			}
		}
	}
}
func DelPkgs(workingDir string, logFile *os.File) {
	allPkgs, errScanDir := os.ReadDir(workingDir)
	if errScanDir != nil {
		logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": Ошибка сканирования директории!\n")
	} else {
		for _, pack := range allPkgs {
			if pack.IsDir() {
				content, errScanDir := os.ReadDir(path.Join(workingDir, pack.Name()))
				if errScanDir != nil {
					logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + ": Ошибка сканирования директории!\n")
				} else {
					if len(content) == 0 {
						os.Remove(path.Join(workingDir, pack.Name()))
					}
				}
			}
		}
	}
}
func (ext SortByExtension) SortFiles(workingDir string, listOfFiles []os.DirEntry) map[string]string {
	pkgs := map[string]string{
		".jpg":        "JPG&JPEGs",
		".jpeg":       "JPG&JPEGs",
		".md":         "MDs&TXTs",
		".txt":        "MDs&TXTs",
		".ppt":        "PPT&PPTXs",
		".pptx":       "PPT&PPTXs",
		".pptx#":      "PPT&PPTXs",
		".ppt#":       "PPT&PPTXs",
		".mp3":        "MP3s",
		".wav":        "WAVs",
		".flac":       "FLACs",
		".epub":       "EPUBs",
		".fb2":        "FB2s",
		".mobi":       "MOBIs",
		".doc":        "DOC&DOCXs",
		".docx":       "DOC&DOCXs",
		".docx#":      "DOC&DOCXs",
		".png":        "PNGs",
		".pdf":        "PDF&TEXs",
		".tex":        "PDF&TEXs",
		".zip":        "ZIP&7Zs",
		".7z":         "ZIP&7Zs",
		".log":        "LOGS SORTING",
		".svg":        "SVGs",
		".rar":        "RARs",
		".xls":        "XLS&XLSXs",
		".xlsx":       "XLS&XLSXs",
		".deb":        "DEBpkgs",
		".flatpakref": "FLATPAKs",
		".flatpak":    "FLATPAKs"}
	return pkgs
}

func (ct SortByContent) SortFiles(workingDir string, listOfFiles []os.DirEntry) map[string]string {
	pkgs := map[string]string{
		".jpg":        "Pictures",
		".png":        "Pictures",
		".gif":        "Pictures",
		".jpeg":       "Pictures",
		".svg":        "Pictures",
		".flac":       "Music",
		".mp3":        "Music",
		".wav":        "Music",
		".aac":        "Music",
		".mp4":        "Videos",
		".avi":        "Videos",
		".mkv":        "Videos",
		".flv":        "Videos",
		".wmv":        "Videos",
		".mov":        "Videos",
		".txt":        "TextDocs",
		".md":         "TextDocs",
		".doc":        "TextDocs",
		".docx":       "TextDocs",
		".tex":        "TextDocs",
		".docx#":      "TextDocs",
		".xls":        "TableDocs",
		".xlsx":       "TableDocs",
		".csv":        "TableDocs",
		".ppt":        "Presentations",
		".pptx":       "Presentations",
		".pdf":        "Presentations",
		".pptx#":      "Presentations",
		".epub":       "Books",
		".fb2":        "Books",
		".mobi":       "Books",
		".atom":       "Books",
		".fb2.zip":    "Books",
		".zip":        "Archives",
		".rar":        "Archives",
		".tar":        "Archives",
		".gz":         "Archives",
		".7z":         "Archives",
		".deb":        "LinuxFiles",
		".flatpakref": "LinuxFiles",
		".flatpak":    "LinuxFiles",
		".run":        "LinuxFiles"}
	return pkgs
}

// -------------------------------------------------------------------------------------------------------------------------
func doSorting(rule GetRules) {
	workingDir, _ := os.UserHomeDir()
	workingDir += "/Загрузки/"
	listOfFiles, _ := os.ReadDir(workingDir)
	pkgs := rule.SortFiles(workingDir, listOfFiles)
	nameLogPack, err := LogPack(pkgs)
	logFile := MakeLogFile(err, workingDir, nameLogPack)
	undoStack, redoStack := Sort(pkgs, workingDir, listOfFiles, logFile)
	DelOldLogs(workingDir, logFile, nameLogPack)
	DelPkgs(workingDir, logFile)
	Undo(&undoStack, &redoStack, logFile)
	Redo(&undoStack, &redoStack, logFile)
	rule.Reminder(workingDir, logFile, pkgs)
	logFile.WriteString("\n")
	defer logFile.Close()
}

// ---------------------------------------------------------------------------------------------------------------------------
func main() {
	byCnt := SortByContent{}
	doSorting(byCnt)
	/*byExt := SortByExtension{}
	doSorting(byExt)*/
}

//3.09 - 4.09 Столкновения имен обрабатывать (давать выбор пользователю) ----> если это возможно, выполнять сравнение и показывать отличия
//4.09 - 9.09 Remind about old files
//remind about doubles
//remind about huge media files
//remind about screenshots
//remind about picture files
//Организация кода (рефакторинг с использованием пакетов)
//Автоматическое отслеживание появления новых файлов для запуска программы
//Рассмотреть возможность создания графинтерфейса на Go
//Вынести настройки в YAML/JSON конфиг
//Ранжирование по времени (папки) внутри папок/

//Следующий крупный шаг - БД и связи
