package logger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"strings"

	"github.com/housinganywhere/hapostman/logger"
)

func getLogs(f string) string {
	content, _ := ioutil.ReadFile(f)
	return string(content)
}

var _ = Describe("HA Postman - Logger", func() {
	Describe("LogLevel consts", func() {
		Context("Constants order", func() {
			It("Fatal < Error", func() {
				Expect(logger.LogFatal).To(BeNumerically("<", logger.LogError))
			})
			It("Error < Warning", func() {
				Expect(logger.LogError).To(BeNumerically("<", logger.LogWarning))
			})
			It("Warning < Info", func() {
				Expect(logger.LogWarning).To(BeNumerically("<", logger.LogInfo))
			})
			It("Info < Debug", func() {
				Expect(logger.LogInfo).To(BeNumerically("<", logger.LogDebug))
			})
			It("Debug < Debug2", func() {
				Expect(logger.LogDebug).To(BeNumerically("<", logger.LogDebug2))
			})
		})
		Context("Constants to String", func() {
			vals := []int{logger.LogFatal, logger.LogError, logger.LogWarning, logger.LogInfo,
				logger.LogDebug, logger.LogDebug2}
			vStrings := []string{"FATAL", "ERROR", "WARNING", "INFO", "DEBUG", "DEBUG2"}
			for i, level := range vals {
				It("level => string", func() {
					Expect(logger.GetLogLevelString(level)).To(Equal(vStrings[i]))
				})
				It("string => level", func() {
					Expect(logger.FindLogLevel(vStrings[i])).To(Equal(level))
				})
			}
		})
		Context("Logger logs only necessary data", func() {
			tmpFile := ""
			BeforeEach(func() {
				tmpfile, err := ioutil.TempFile("", "tests")
				if err != nil {
					panic(err)
				}
				tmpFile = tmpfile.Name()
			})
			AfterEach(func() {
				if tmpFile != "" {
					os.Remove(tmpFile)
					tmpFile = ""
				}
			})
			It("level=Fatal (fatal called)", func() {
				defer func() {
					r := recover()
					Expect(r).NotTo(BeNil())
					content := getLogs(tmpFile)
					Expect(content).To(ContainSubstring("FATAL TestF"))
					Expect(len(strings.Split(content, "\n"))).To(Equal(2))
				}()
				l := logger.NewLogger(logger.LogFatal, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				l.Fatal("TestF")
			})
			It("level=Fatal (no  Fatal called)", func() {
				l := logger.NewLogger(logger.LogFatal, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				content := getLogs(tmpFile)
				Expect(content).To(Equal(""))
			})
			It("level=Error", func() {
				l := logger.NewLogger(logger.LogError, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				content := getLogs(tmpFile)
				Expect(content).To(ContainSubstring("ERROR TestE"))
				Expect(content).NotTo(ContainSubstring("Warning TestW"))
				Expect(content).NotTo(ContainSubstring("INFO TestI"))
				Expect(content).NotTo(ContainSubstring("DEBUG TestD"))
				Expect(content).NotTo(ContainSubstring("DEBUG2 TestD2"))
				Expect(len(strings.Split(content, "\n"))).To(Equal(2))
			})
			It("level=Warning", func() {
				l := logger.NewLogger(logger.LogWarning, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				content := getLogs(tmpFile)
				Expect(content).To(ContainSubstring("ERROR TestE"))
				Expect(content).To(ContainSubstring("WARNING TestW"))
				Expect(content).NotTo(ContainSubstring("INFO TestI"))
				Expect(content).NotTo(ContainSubstring("DEBUG TestD"))
				Expect(content).NotTo(ContainSubstring("DEBUG2 TestD2"))
				Expect(len(strings.Split(content, "\n"))).To(Equal(3))
			})
			It("level=Info", func() {
				l := logger.NewLogger(logger.LogInfo, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				content := getLogs(tmpFile)
				Expect(content).To(ContainSubstring("ERROR TestE"))
				Expect(content).To(ContainSubstring("WARNING TestW"))
				Expect(content).To(ContainSubstring("INFO TestI"))
				Expect(content).NotTo(ContainSubstring("DEBUG TestD"))
				Expect(content).NotTo(ContainSubstring("DEBUG2 TestD2"))
				Expect(len(strings.Split(content, "\n"))).To(Equal(4))
			})
			It("level=Debug", func() {
				l := logger.NewLogger(logger.LogDebug, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				content := getLogs(tmpFile)
				Expect(content).To(ContainSubstring("ERROR TestE"))
				Expect(content).To(ContainSubstring("WARNING TestW"))
				Expect(content).To(ContainSubstring("INFO TestI"))
				Expect(content).To(ContainSubstring("DEBUG TestD"))
				Expect(content).NotTo(ContainSubstring("DEBUG2 TestD2"))
				Expect(len(strings.Split(content, "\n"))).To(Equal(5))
			})
			It("level=Debug2", func() {
				l := logger.NewLogger(logger.LogDebug2, tmpFile)
				l.Debug2("TestD2")
				l.Debug("TestD")
				l.Info("TestI")
				l.Warning("TestW")
				l.Error("TestE")
				content := getLogs(tmpFile)
				Expect(content).To(ContainSubstring("ERROR TestE"))
				Expect(content).To(ContainSubstring("WARNING TestW"))
				Expect(content).To(ContainSubstring("INFO TestI"))
				Expect(content).To(ContainSubstring("DEBUG TestD"))
				Expect(content).To(ContainSubstring("DEBUG2 TestD2"))
				Expect(len(strings.Split(content, "\n"))).To(Equal(6))
			})
		})
	})
})
