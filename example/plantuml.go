package main

//import (
//
//	"bytes"
//	"fmt"
//	"io"
//	"os"
//	"os/exec"
//)
//
func main() {
}

//}
//	PlantUML(`@startuml
//Bob -> Alice : hello
//@enduml`)
//}
//func PlantUML(plantumlString string) string {
//	//create command
//	echo := exec.Command("echo", `"`+plantumlString+`"`)
//	plantuml := exec.Command("java", "-Djava.awt.headless=true", "-jar", "plantuml.jar", "-p", "-tpng")
//
//	//make a pipe
//	reader, writer := io.Pipe()
//	var buf bytes.Buffer
//
//	//set the output of "cat" command to pipe writer
//	echo.Stdout = writer
//	//set the input of the "wc" command pipe reader
//
//	plantuml.Stdin = reader
//
//	//cache the output of "wc" to memory
//	plantuml.Stdout = &buf
//
//	//start to execute "cat" command
//	echo.Start()
//
//	//start to execute "wc" command
//	plantuml.Start()
//
//	//waiting for "cat" command complete and close the writer
//	echo.Wait()
//	writer.Close()
//
//	//waiting for the "wc" command complete and close the reader
//	plantuml.Wait()
//	reader.Close()
//	//copy the buf to the standard output
//	//io.Copy(os.Stdout, &buf)
//	fmt.Println(buf.String())
//	f, err := os.Create("/out/test.png")
//	if err != nil {
//		panic(err)
//	}
//	_, err = f.Write(buf.Bytes())
//	if err != nil {
//		panic(err)
//	}
//}
