package barcode

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// A device is an object representing a barcode device.
type device struct {
	deviceFile *os.File      // Open file object for the serial device
	buffer     *bytes.Buffer // Buffer to store the characters read from the device
	filename   string        // File path of the serial device
}

// Creates and initializes a device
func NewDevice(filename string) (this *device, err error) {
	//Open serial device file
	dev_file, err := os.OpenFile(filename, os.O_RDONLY, 0)
	//TODO: add something to close files when stopped
	if err != nil {
		return nil, fmt.Errorf("Error opening barcode device file='%s' error='%s'", filename, err)
	}

	//Create device
	this = &device{
		buffer:     &bytes.Buffer{},
		filename:   filename,
		deviceFile: dev_file,
	}

	return this, nil
}

// Read a chunk of messages from the device
func (this *device) syncMessages() (err error) {

	// make temporary buffer
	buf := make([]byte, 1024)

	// Read from serial device
	_, err = this.deviceFile.Read(buf)
	if err != nil {
		return fmt.Errorf("Error syncing barcode device messages: %s", err)
	}

	//trim null bytes
	buf = bytes.Trim(buf, "\x00")
	s_buf := string(buf)
	s_buf = strings.Replace(s_buf, "\r", "\n", -1)

	// write message to real buffer
	this.buffer.WriteString(s_buf)
	//TODO add locking here

	return nil
}

// peekMessages checks for new complete messages in the buffer.
func (this *device) peekMessages() (messages []string, partial string) {
	// separate buffer into messages
	messages = strings.Split(this.buffer.String(), "\n")

	// get last message in buffer
	lastmessage := messages[len(messages)-1]
	messages = messages[:len(messages)-1]

	// if the message isn't complete, return it as partial
	//TODO check if this logic is still right. it was changed.
	if lastmessage != "" {
		partial = lastmessage
	}

	return messages, partial
}

// PopNewMessages removes any new messages from the buffer and returns them.
func (this *device) popNewMessages() (messages []string) {

	// Get new messages
	messages, partial := this.peekMessages()

	// Remove the messages from the buffer
	//TODO check if this logic is still right.
	if partial != "" {
		// If there is an incomplete message, reinitialize
		// the buffer with just the incomplete one
		this.buffer = bytes.NewBufferString(partial)
	} else {
		// If we don't have any partial messages, then just make a new empty buffer
		this.buffer = &bytes.Buffer{}
	}

	return messages
}

func (this *device) Run() (chan string, chan error) {
	out := make(chan string)
	err := make(chan error)
	go this.run(out, err)
	return out, err
}

func (this *device) run(outchan chan string, errchan chan error) {
	//begin loop
	for {
		//sync messages with barcode device

		err := this.syncMessages()
		if err != nil {

			errchan <- err
		}

		messages := this.popNewMessages()

		for _, message := range messages {
			// make barcode event
			outchan <- message
		}
	}
}
