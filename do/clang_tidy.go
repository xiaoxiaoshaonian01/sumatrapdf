package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*
https://clang.llvm.org/extra/clang-tidy/checks/list.html
https://codeyarns.com/2019/01/28/how-to-use-clang-tidy/
https://www.reddit.com/r/cpp/comments/ezn21f/which_checks_do_you_use_for_clangtidy/
https://www.reddit.com/r/cpp/comments/5bqkk5/good_clangtidy_files/
https://www.reddit.com/r/cpp/comments/7obg9p/how_do_you_use_clangtidy/
https://github.com/KratosMultiphysics/Kratos/wiki/How-to-use-Clang-Tidy-to-automatically-correct-code
https://sarcasm.github.io/notes/dev/clang-tidy.html
https://www.labri.fr/perso/fleury/posts/programming/using-clang-tidy-and-clang-format.html
*/

/*
ad-hoc execution:
clang-tidy.exe --checks=-clang-diagnostic-microsoft-goto,-clang-diagnostic-unused-value -extra-arg=-std=c++20 .\src\*.cpp -- -I mupdf/include -I src -I src/utils -I src/wingui -I ext/WDL -DUNICODE -DWIN32 -D_WIN32 -D_CRT_SECURE_NO_WARNINGS -DWINVER=0x0a00 -D_WIN32_WINNT=0x0a00

ls src\utils\*.cpp | select Name

clang-tidy src/utils/ColorUtil.cpp -fix -checks="-*,readability-braces-around-statements" -extra-arg=-std=c++20 -- -I mupdf/include -I src -I src/utils -I src/wingui -I ext/WDL -I ext/unarr -DUNICODE -DWIN32 -D_WIN32 -D_CRT_SECURE_NO_WARNINGS -DWINVER=0x0a00 -D_WIN32_WINNT=0x0a00
*/

/*
src\utils:

CryptoUtil.cpp
CssParser.cpp
DbgHelpDyn.cpp
Dict.cpp
DirIter.cpp
Dpi.cpp
FileUtil.cpp
FileWatcher.cpp
FzImgReader.cpp
GdiPlusUtil.cpp
GuessFileType.cpp
HtmlParserLookup.cpp
HtmlPrettyPrint.cpp
HtmlPullParser.cpp
HtmlWindow.cpp
HttpUtil.cpp
JsonParser.cpp
Log.cpp
LogDbg.cpp
LzmaSimpleArchive.cpp
PalmDbReader.cpp
RegistryPaths.cpp
SerializeTxt.cpp
SettingsUtil.cpp
SquareTreeParser.cpp
StrconvUtil.cpp
StrFormat.cpp
StringViewUtil.cpp
StrSlice.cpp
StrUtil.cpp
StrUtil_win.cpp
TgaReader.cpp
ThreadUtil.cpp
TrivialHtmlParser.cpp
TxtParser.cpp
UITask.cpp
UtAssert.cpp
WebpReader.cpp
WinDynCalls.cpp
WinUtil.cpp
ZipUtil.cpp
*/

const clangTidyLogFile = "clangtidy.out.txt"

func clangTidyFile(path string) {
	args := []string{
		"--checks=-clang-diagnostic-microsoft-goto,-clang-diagnostic-unused-value",
		"-extra-arg=-std=c++20",
		"", // file
		"--",
		"-I", "mupdf/include",
		"-I", "src",
		"-I", "src/utils",
		"-I", "src/wingui",
		"-I", "ext/WDL",
		"-I", "ext/CHMLib",
		"-I", "ext/libdjvu",
		"-I", "ext/zlib",
		"-I", "ext/synctex",
		"-I", "ext/unarr",
		"-I", "ext/lzma/C",
		"-I", "ext/libwebp/src",
		"-I", "ext/freetype/include",

		"-DUNICODE",
		"-DWIN32",
		"-D_WIN32",
		"-D_CRT_SECURE_NO_WARNINGS",
		"-DWINVER=0x0a00",
		"-D_WIN32_WINNT=0x0a00",
		"-DPRE_RELEASE_VER=3.3",
	}
	args[2] = path
	cmd := exec.Command("clang-tidy", args...)
	_ = runCmdShowProgressAndLog(cmd, clangTidyLogFile)
}

func runClangTidy() {
	os.Remove(clangTidyLogFile)
	files := []string{
		`src\*.cpp`,
		`src\*.h`,
		`src\mui\*.cpp`,
		`src\mui\*.h`,
		`src\utils\*.cpp`,
		`src\utils\*.h`,
		`src\utils\tests\*.cpp`,
		`src\utils\tests\*.h`,
		`src\wingui\*.cpp`,
		`src\wingui\*.h`,
		`src\uia\*.cpp`,
		`src\uia\*.h`,
		`src\tools\*.cpp`,
		`src\tools\*.h`,
		`ext\mupdf_load_system_font.c`,
	}

	isWhiteListed := func(s string) bool {
		whitelisted := []string{
			"resource.h",
			"Version.h",
			"Trans_sumatra_txt.cpp",
			"Trans_installer_txt.cpp",
		}
		s = strings.ToLower(s)

		if strings.HasSuffix(s, ".h") {
			return true
		}

		for _, wl := range whitelisted {
			wl = strings.ToLower(wl)
			if strings.Contains(s, wl) {
				logf("Whitelisted '%s'\n", s)
				return true
			}
		}
		return false
	}
	for _, globPattern := range files {
		paths, err := filepath.Glob(globPattern)
		must(err)
		for _, path := range paths {
			if isWhiteListed(path) {
				continue
			}
			clangTidyFile(path)
		}
	}
	logf("\nLogged output to '%s'\n", clangTidyLogFile)
}