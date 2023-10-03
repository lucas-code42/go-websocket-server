function textMsg() {
    const ws = new WebSocket("ws://localhost:8080/ws-text");
    ws.send("hello from web-socket");
}
textMsg();

function jsonMsg() {
    const ws = new WebSocket("ws://localhost:8080/ws-json");
    ws.send(
        JSON.stringify({
            protocol: "ws",
            msg: "hello from ws client",
            status: "open"
        })
    );
    ws.onmessage = function (event) { console.log(event.data) }
}
jsonMsg();