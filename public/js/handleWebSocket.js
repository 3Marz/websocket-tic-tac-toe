
const loading = document.getElementById("loading")
const message = document.getElementById("message")

const connectButton = document.getElementById('connect');
const closeButton = document.getElementById('close');
var ws;
var yourTurn = false;

function toggleLoading() {
	loading.classList.toggle("hide-down")
}

connectButton.addEventListener('click', () => {
	if(ws) {return}
	ws = new WebSocket("http://localhost:5050/ws");
	ws.onopen = function(evt) {
		connectButton.innerText = 'Connected';
		loading.innerText = "Wating For An Opponent..."
		console.log('Connected');
	}
	ws.onclose = function(evt) {
		resetEverything()
	}
	ws.onmessage = function(evt) {
		//console.log("RESPONSE: " + evt.data);
		//console.log(evt)
		switch (evt.data) {
			case "Got A Game":
				loading.innerText = "Found An Opponent"
				toggleLoading()
				message.innerText = "Opponent Turn..."
				break
			case "Your Turn":
				yourTurn = true
				message.innerText = "Your Turn"
				break
			case "You Won":
				loading.innerText = "You Won"
				if(ws) { ws.close() }
				break
			case "You Lost":
				loading.innerText = "You Lost"
				if(ws) { ws.close() }
				break
			case "Opponent left the game":
				loading.innerText = "Opponent left the game"
				if(ws) {
					ws.close()
				}
				break
			default:
				let board = JSON.parse(evt.data)
				writeBoard(board)
				break
		}
	}
	ws.onerror = function(evt) {
		console.log("ERROR: " + evt.data);
	}
})

closeButton.addEventListener('click', () => {
	if (!ws) {
		return
	}
	ws.close();
})

function resetBoard() {
	for (const span of document.querySelectorAll('span')) {
		span.innerText = ''
	}
}
function readBoard() {
	let board = [[], [], []];
	for (const span of document.querySelectorAll('span')) {
		board[span.getAttribute('row')][span.getAttribute('col')] = span.innerText
	}
	return board
}
function writeBoard(board) {
	for (const span of document.querySelectorAll('span')) {
		span.innerText = board[span.getAttribute('row')][span.getAttribute('col')]
	}
}

function resetEverything() {
	connectButton.innerText = 'Connect';
	console.log('Connection closed');
	ws = null
	loading.innerText = loading.innerText == "Found An Opponent" ? "Left The Game" : loading.innerText
	message.innerText = "..."
	toggleLoading()
	resetBoard()
}

function sendPlay(element) {
	if (!ws) {
		return
	}
	if (!yourTurn) {
		return
	}
	if(element.innerText !== '') {
		return
	}
	//let board = readBoard()
	//board[element.getAttribute('row')][element.getAttribute('col')] = 'X'
	//console.log(JSON.stringify(board))
	let req = [
		Number(element.getAttribute('row')),
		Number(element.getAttribute('col'))
	]
	console.log(JSON.stringify(req))
	ws.send(JSON.stringify(req))
	yourTurn = false
	message.innerText = "Opponent Turn..."
}

