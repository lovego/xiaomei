package cmd

import (
	"io"
	"net/http"
	"os/exec"
)

// `bash`, `-c`, `for ((i=0; i<10; i++)); do echo $i; sleep 1; done`
func Stream(stream io.Writer, name string, args ...string) {
	cmd := exec.Command(name, args...)
	flusher, ok := stream.(http.Flusher)
	if !ok {
		cmd.Stdout = stream
		cmd.Stderr = stream
		cmd.Run()
		return
	}

	if rw, ok := stream.(http.ResponseWriter); ok {
		// rw.Header().Set("Connection", "Keep-Alive")
		// rw.Header().Set("Transfer-Encoding", "chunked")
		rw.Header().Set("X-Content-Type-Options", "nosniff")
	}

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	cmd.Stdout = writer
	cmd.Stderr = writer

	go func() {
		cmd.Run()
		writer.Close()
	}()

	loopReadWrite(reader, stream, flusher)
}

func loopReadWrite(reader *io.PipeReader, stream io.Writer, flusher http.Flusher) {
	buffer := make([]byte, 100)
	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			stream.Write(buffer[:n])
			flusher.Flush()
		}
		if err != nil {
			break
		}
	}
}
