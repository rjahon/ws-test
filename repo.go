package main

type APIClient interface {
	/*
		the function should implement a connection to the exchange's WebSocket.
		If there is a connection problem, it should return an error
	*/
	Connection() error

	/*
		Implement a disconnect function from websocket
	*/
	Disconnect()

	/*
		Implement a function that will subscribe to updates
		of BBO for a given symbol

		The symbol must be of the form "TOKEN_ASSET"
		As an example "USDT_BTC" where USDT is TOKEN and BTC is ASSET

		You will need to convert the symbol in such a way that
		it complies with the exchange standard.

		If there are problems with listening, it should return an error and correctly disconnect from WebSocket.
	*/
	SubscribeToChannel(symbol string) error

	/*
		Implement a function that will write the data that
		we receive from the exchange websocket to the channel
	*/
	ReadMessagesFromChannel(ch chan<- BestOrderBook)

	/*
		Implement a function that will support connecting to a websocket
	*/
	WriteMessagesToChannel()
}
