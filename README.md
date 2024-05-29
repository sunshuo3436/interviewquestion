# TCP Server for Parsing Custom Protocol

This is a TCP server written in Go that listens for incoming connections, parses a custom protocol, and prints the parsed data. The custom protocol defines a specific message structure that includes a Device ID encoded in ASCII.

## License

This project is licensed under the MIT License.
## Prerequisites

Before running the server, ensure you have the following installed:

- Go (version 1.16 or higher)

## Installation

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/server.git
```
## Usage

1.Navigate to the project directory:
```bash
cd server
```
2.Run the server:
```bash
go run server.go
```
3.In a separate terminal window, run the client to send test messages:
bash
```bash
go run client.go
```
## Custom Protocol

The custom protocol used by this TCP server has a specific message structure, as described below:

- **Start Bytes**: 2 bytes (fixed value `0xaa 0xaa`)
- **Length**: 2 bytes (length of the data unit)
- **AFN**: 1 byte (application layer function code)
- **Device ID**: 10 bytes (ASCII encoded)
- **Socket Number**: 1 byte (number indicating socket)
- **Data Unit**: Variable length (actual data)

A sample message adhering to this protocol structure looks like this:
```bash
aa aa 00 1b 04 32 30 32 32 30 30 30 33 32 39 02 17 09 0d 0a 0c 01 01 00 00 00 64 00 55 55 55 55 55 6a
```

In this example:
- Start Bytes: `0xaa 0xaa`
- Length: `0x1b` (27 in decimal)
- AFN: `0x04`
- Device ID: ASCII encoded string `2022000329`
- Socket Number: `0x02`
- Data Unit: `17090D0A0C010100000064005555`

## My Challenges and Solutions

### 1. Limited Understanding of TCP Communication

My solution was to learn and understand TCP communication principles by studying other people's code and mimicking them with the help of GPT. I will continue to delve into relevant documentation and resources to enhance my understanding and skills.

### 2. Mistaken Use of Postman, Switched to Testing with client.go File

Due to the mistake in using Postman, I opted to test server.go using the client.go file instead. Nevertheless, I intend to retry using Postman for testing when time permits, to comprehensively evaluate the functionality and performance of the server.

### 3. Ensuring Code Readability, Especially with ASCII Encoded Device ID

I've recognized the significance of code readability for comprehension and maintenance, particularly with fields like Device ID, which utilize ASCII encoding. To ensure code clarity, I'll incorporate appropriate comments and strive for clear and concise code writing.




# Third Revision

## Handling Sticky and Fragmented Messages

### Key Points:
1. **Setting `dataUnitStart` and `dataUnitEnd` Calculation**: Ensure the `DataUnit` start and end positions are calculated correctly to avoid extra data.
2. **Boundary Check**: Ensure `dataUnitEnd` does not exceed the actual message length to prevent index out of range errors.

### `printMessage` Function
The `printMessage` function is responsible for extracting and printing the different parts of the message. Here's the updated function:

```go
func printMessage(message []byte) {
    // Extract and print different parts of the message
    startBytes := message[:StartBytesLen]
    lengthBytes := message[StartBytesLen : StartBytesLen+LengthBytesLen]
    length := binary.BigEndian.Uint16(lengthBytes)
    afn := message[StartBytesLen+LengthBytesLen]

    // Calculate the start and end positions of DeviceID
    deviceIDStart := StartBytesLen + LengthBytesLen + AFNBytesLen
    deviceIDEnd := deviceIDStart + DeviceIDLen
    deviceID := message[deviceIDStart:deviceIDEnd]

    socketNum := message[deviceIDEnd]

    // Calculate the start position and length of DataUnit
    dataUnitStart := HeaderLen
    dataUnitEnd := HeaderLen + int(length) - 1
    if dataUnitEnd > len(message) {
        dataUnitEnd = len(message)
    }
    dataUnit := message[dataUnitStart:dataUnitEnd]

    fmt.Printf("Start Bytes: %X\n", startBytes)
    fmt.Printf("Length: %d\n", length)
    fmt.Printf("AFN: %X\n", afn)
    fmt.Printf("Device ID: %s\n", string(deviceID))
    fmt.Printf("Socket Number: %d\n", socketNum)
    fmt.Printf("Data Unit: %X\n", dataUnit)
}

