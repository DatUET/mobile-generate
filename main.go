package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type FvmConfig struct {
	FlutterSdkVersion string `json:"flutterSdkVersion"`
}

type Config struct {
	ProjectName    string `json:"projectName"`
	PackageName    string `json:"packageName"`
	AppName        string `json:"appName"`
	Pattern        string `json:"pattern"`
	FlutterVersion string `json:"flutterVersion"`
	CreateAt       string `json:"createAt"`
}

const Version = "0.1.0"

func fileExists() bool {
	_, err := os.Stat(".dmg/config.json")
	return !os.IsNotExist(err) // Trả về true nếu file tồn tại
}

func getParttern() string {
	// Mở file
	filePath := ".dmg/config.json" // Đường dẫn file JSON
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	// Decode file JSON
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}
	return config.Pattern
}

func getFvmFlutterVersion() string {
	// Mở file
	filePath := ".fvm/fvm_config.json" // Đường dẫn file JSON
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	// Decode file JSON
	var config FvmConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return ""
	}
	return config.FlutterSdkVersion
}

func printTutorial() {
	fmt.Println(`Usage: dmg [OPTION]
OPTION:
  create        Create new Flutter project
  new_page      Create new screen
  new_stack     Create new stack
  build			Build app
  pub_get       Run pub get
  setup         Set up project
  --version     Check version
  -h, --help    Helper`)
}

func printTutorialNewPage() {
	fmt.Println(`Usage: dmg new_page [OPTION]
OPTION:
  -n            Name page (Required)
  -p            Prefix name file (Required)
  -t            Screen type, default is normal screen (Optional)
EXAMPLE:
  dmg new_page -n Home -p home -t list`)
}

func printTutorialNewStack() {
	fmt.Println(`Usage: dmg new_stack [OPTION]
OPTION:
  -n            	Name stack (Required)
  -p            	Prefix name file (Required)
  --enable-local	Generate local data source (Optional)
EXAMPLE:
  dmg new_stack -n User -p user --enable-local`)
}

func printTutorialBuild() {
	fmt.Println("┌────────────────────────────────────────────────────────┐")
	fmt.Println("│                     Auto Build App                     │")
	fmt.Println("├────────────────────────────────────────────────────────┤")
	fmt.Println("│       COMMAND: dmg build [Environment] [Tester]        │")
	fmt.Println("├───────────────┬────────────────────────────────────────┤")
	fmt.Println("│               │1.  --android                           │")
	fmt.Println("│               │2.  --dev-debug-android                 │")
	fmt.Println("│               │3.  --staging-staging-android           │")
	fmt.Println("│               │4.  --prod-release-android              │")
	fmt.Println("│               │5.  --store-android                     │")
	fmt.Println("│               │6.  --ios                               │")
	fmt.Println("│ [Environment] │7.  --dev-debug-ios                     │")
	fmt.Println("│               │8.  --staging-staging-ios               │")
	fmt.Println("│               │9.  --prod-release-ios                  │")
	fmt.Println("│               │10. --store-ios                         │")
	fmt.Println("│               │11. --dev-debug                         │")
	fmt.Println("│               │12. --staging-staging                   │")
	fmt.Println("│               │13. --prod-release                      │")
	fmt.Println("│               │14. --store                             │")
	fmt.Println("├───────────────┼────────────────────────────────────────┤")
	fmt.Println("│               │1.  --dev                               │")
	fmt.Println("│   [Tester]    │2.  --tester                            │")
	fmt.Println("│               │3.  --client                            │")
	fmt.Println("├───────────────┴────────────────────────────────────────┤")
	fmt.Println("│                                                        │")
	fmt.Println("│Or Please read README to learn more.                    │")
	fmt.Println("└────────────────────────────────────────────────────────┘")
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func runNewPage(namePage string, prefix string, screenType string) {
	if !fileExists() {
		fmt.Println("Error: No such dmg config file")
		os.Exit(1)
	}

	if namePage == "" || prefix == "" {
		printTutorialNewPage()
		os.Exit(1)
	}

	pattern := getParttern()
	switch pattern {
	case "BLoC":
		runCommand("new_page.sh", "-n", namePage, "-p", prefix, "-t", screenType)
	case "GetX":
		runCommand("new_page_getx.sh", "-n", namePage, "-p", prefix, "-t", screenType)
	case "Riverpod":
		fmt.Println("Currently, automatic screen generation for Riverpod pattern is not supported. We will update later.")
	default:
		fmt.Println("Invalid pattern")
	}
}

func runNewStack(namePage string, prefix string, enableLocal bool) {
	if !fileExists() {
		fmt.Println("Error: No such dmg config file")
		os.Exit(1)
	}

	if namePage == "" || prefix == "" {
		printTutorialNewPage()
		os.Exit(1)
	}

	pattern := getParttern()
	switch pattern {
	case "BLoC":
		if enableLocal {
			runCommand("new_stack.sh", "-n", namePage, "-p", prefix, "--enable-local")
		} else {
			runCommand("new_stack.sh", "-n", namePage, "-p", prefix)
		}
	case "GetX":
		if enableLocal {
			runCommand("new_stack_getx.sh", "-n", namePage, "-p", prefix, "--enable-local")
		} else {
			runCommand("new_stack_getx.sh", "-n", namePage, "-p", prefix)
		}
	case "Riverpod":
		if enableLocal {
			runCommand("new_stack.sh", "-n", namePage, "-p", prefix, "--enable-local")
		} else {
			runCommand("new_stack.sh", "-n", namePage, "-p", prefix)
		}
	default:
		fmt.Println("Invalid pattern")
	}
}

func runBuildOption(args []string) {
	listBuildType := []string{
		"--android",
		"--dev-debug-android",
		"--staging-staging-android",
		"--prod-release-android",
		"--store-android",
		"--ios",
		"--dev-debug-ios",
		"--staging-staging-ios",
		"--prod-release-ios",
		"--store-ios",
		"--dev-debug",
		"--staging-staging",
		"--prod-release",
		"--store",
	}

	listTestter := []string{
		"--dev",
		"--tester",
		"--client",
	}
	if len(os.Args) < 3 {
		printTutorialBuild()
		os.Exit(1)
	}
	buildType := args[2]
	if buildType != "store-ios" && buildType != "store-android" && buildType != "store" && len(args) < 4 {
		printTutorialBuild()
		os.Exit(1)
	}
	tester := ""
	if len(args) == 4 {
		tester = args[3]
	}
	if !contains(listBuildType, buildType) || !contains(listTestter, tester) {
		printTutorialBuild()
		os.Exit(1)
	}
	runCommand("build.sh", buildType, tester)
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

func updateFvmLinked() {
	if !fileExists() {
		fmt.Println("Error: No such dmg config file")
		os.Exit(1)
	}
	fvnFlutterVersion := getFvmFlutterVersion()
	filePath := ".dmg/config.json" // Đường dẫn file JSON

	// Đọc nội dung file JSON
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse JSON thành struct
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Cập nhật giá trị flutterVersion
	config.FlutterVersion = fvnFlutterVersion // Giá trị mới

	// Chuyển struct về JSON
	updatedJSON, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Ghi lại vào file
	err = os.WriteFile(filePath, updatedJSON, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

func main() {
	updateFvmLinked()
	if len(os.Args) < 2 {
		printTutorial()
		os.Exit(1)
	}

	runOption := os.Args[1]

	// Define flags
	namePage := flag.String("n", "", "Name of the page/stack")
	prefix := flag.String("p", "", "Prefix for naming")
	screenType := flag.String("t", "", "Type of screen")
	enableLocal := flag.Bool("enable-local", false, "Enable local data source")
	showHelp := flag.Bool("h", false, "Show help")
	showHelpLong := flag.Bool("help", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")

	if runOption != "build" {
		err := flag.CommandLine.Parse(os.Args[2:])
		if err != nil {
			printTutorial()
			os.Exit(1)
		}
	}

	// Handle --help and --version
	if *showHelp || *showHelpLong {
		printTutorial()
		os.Exit(0)
	}
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	// Handle different options
	switch runOption {
	case "create":
		runCommand("sli.sh")
	case "new_page":
		runNewPage(*namePage, *prefix, *screenType)
	case "new_stack":
		runNewStack(*namePage, *prefix, *enableLocal)
	case "build":
		if fileExists() {
			runBuildOption(os.Args)
		} else {
			fmt.Println("Error: No such dmg config file")
			os.Exit(1)
		}
	case "pub_get":
		if fileExists() {
			runCommand("pub_get.sh")
		} else {
			fmt.Println("Error: No such dmg config file")
			os.Exit(1)
		}
	case "setup":
		if fileExists() {
			runCommand("setup.sh")
		} else {
			fmt.Println("Error: No such dmg config file")
			os.Exit(1)
		}
	case "-h":
		printTutorial()
		os.Exit(1)
	case "--help":
		printTutorial()
		os.Exit(1)
	default:
		fmt.Println("Invalid option:", runOption)
		printTutorial()
		os.Exit(1)
	}
}
