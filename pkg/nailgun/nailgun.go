package nailgun

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type NailgunConnection struct {
	conn net.Conn
}

func (self *NailgunConnection) sendChunk(chunkType byte, payload []byte) (err error) {
	var header [5]byte

	n := len(payload)
	header[0] = byte((n >> 24) & 0xff)
	header[1] = byte((n >> 16) & 0xff)
	header[2] = byte((n >> 8) & 0xff)
	header[3] = byte((n >> 0) & 0xff)

	header[4] = chunkType

	n, err = self.conn.Write(header[:])

	if err != nil {
		return
	}

	if n != len(header) {
		log.Panic("Unexpected short-write")
	}

	n, err = self.conn.Write(payload)
	if err != nil {
		return
	}

	if n != len(payload) {
		log.Panic("Unexpected short-write")
	}

	return nil
}

func (self *NailgunConnection) sendArguments(args []string) (err error) {
	for i, arg := range args {
		if i <= 1 {
			// Unlike nailgun, we always skip the first two args
			// the first is our name
			// the second is the command to run
			continue
		}

		err = self.sendChunk('A', []byte(arg))
		if err != nil {
			return
		}
	}

	return nil
}

func (self *NailgunConnection) sendEnvironment(env []string) (err error) {
	for _, env := range env {
		err = self.sendChunk('E', []byte(env))
		if err != nil {
			return
		}
	}

	return nil
}

func (self *NailgunConnection) sendWorkingDirectory() (err error) {
	cwd, err := filepath.Abs(".")
	if err != nil {
		return
	}

	err = self.sendChunk('D', []byte(cwd))
	if err != nil {
		return
	}

	return nil
}

func (self *NailgunConnection) sendCommand(args []string) (err error) {
	command := args[0]

	err = self.sendChunk('C', []byte(command))
	if err != nil {
		return
	}

	return nil
}

func (self *NailgunConnection) forwardStdin(r io.Reader) (err error) {
	N := 8192
	buffer := make([]byte, N)
	var n int

	for {
		n, err = r.Read(buffer)
		if err != nil {
			if n == 0 && err == io.EOF {
				break
			}
			return
		}

		err = self.sendChunk('0', buffer[:n])
		if err != nil {
			return
		}
	}

	err = self.sendChunk('.', buffer[0:0])
	if err != nil {
		return
	}

	return nil
}

func (self *NailgunConnection) readFully(buffer []byte) (err error) {
	var n int
	n, err = io.ReadFull(self.conn, buffer)
	if err != nil {
		return
	}

	if n != len(buffer) {
		log.Panic("Unexpected short read")
	}

	return nil
}

func (self *NailgunConnection) readFromServer(dest io.Writer) (exitCode int, err error) {
	N := 8192
	header := make([]byte, 5, 5)
	buffer := make([]byte, N, N)

	var n int

	for {
		err = self.readFully(header[0:5])
		if err != nil {
			return
		}

		payloadLength := (int(header[0]) << 24)
		payloadLength |= (int(header[1]) << 16)
		payloadLength |= (int(header[2]) << 8)
		payloadLength |= (int(header[3]) << 0)

		payloadType := header[4]

		//log.Printf("Chunk %v %v %v\n", header, payloadLength, payloadType)

		switch payloadType {
		case '1':
			dest = os.Stdout
		case '2':
			dest = os.Stderr

		case 'S':
			// Undocumented chunk type "start input"
			if payloadLength != 0 {
				err = fmt.Errorf("Expected 0 length for S chunk")
				return
			}
			continue

		case 'X':
			err = self.readFully(buffer[0:payloadLength])
			if err != nil {
				return
			}

			exitCode, err = strconv.Atoi(strings.ReplaceAll(string(buffer[0:payloadLength]), "\n", ""))
			if err != nil {
				return
			}

			return exitCode, nil

		default:
			err = fmt.Errorf("Unexpected chunk type %v", payloadType)
			return
		}

		for payloadLength > 0 {
			read := payloadLength
			if read > N {
				read = N
			}
			n, err = self.conn.Read(buffer[0:read])
			if n > 0 {
				_, writeErr := dest.Write(buffer[0:n])
				if writeErr != nil {
					return 0, writeErr
				}

				payloadLength -= n
			}

			if err != nil {
				return
			}

		}
	}

	log.Panic("Unreachable")

	return
}

func (self *NailgunConnection) close() (err error) {
	err = self.conn.Close()
	if err != nil {
		return
	}
	return nil
}

func findEnvironmentVariable(key string) (value string) {
	prefix := key + "="
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix) {
			value = env[len(prefix):]
			return
		}
	}
	return ""
}

func runCommand(argv []string) {
	log.Printf("Starting server: %v\n", argv)
	cmd := exec.Command(argv[0], argv[1:]...)
	//cmd.Stdin = strings.NewReader("some input")
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Run(args []string, env []string, stdin io.Reader, stdout io.Writer) {
	// if len(os.Args) <= 2 {
	// 	log.Printf("Must pass (at least) command/class to run")
	// 	os.Exit(2)
	// }

	port := findEnvironmentVariable("NAILGUN_PORT")
	if port == "" {
		port = "2113"
	}

	bootstrap := args[1] == "bootstrap"

	var err error
	var conn net.Conn

	if bootstrap {
		conn, err = net.Dial("tcp", "127.0.0.1:"+port)
		if err != nil {
			//			os.Stdin.Close()
			//			os.Stdout.Close()
			//			os.Stderr.Close()
			runCommand(args[2:])
		} else {
			conn.Close()
		}
		os.Exit(0)
	}

	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", "127.0.0.1:"+port)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}

	if err != nil {
		log.Printf("Unable to connect to nailgun server: %v\n", err)
		os.Exit(1)
	}

	ng := &NailgunConnection{}
	ng.conn = conn

	ng.sendArguments(args)
	ng.sendEnvironment(env)
	ng.sendWorkingDirectory()
	ng.sendCommand(args)

	go func() {
		err := ng.forwardStdin(stdin)
		if err != nil {
			log.Printf("Error forwarding stdin: %v\n", err)
			ng.close()
		}
	}()

	exitCode, err := ng.readFromServer(stdout)
	ng.close()

	if err != nil {
		if err != io.EOF {
			log.Printf("Error communicating with background process: %v\n", err)
			os.Exit(2)
		}
	}

	os.Exit(exitCode)
}
