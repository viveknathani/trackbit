const socket = new WebSocket('ws://localhost:8081/stats');

socket.onopen = function(event) {
    console.log('connection established.');
    socket.send(JSON.stringify({
        Name: "eth0"
    }));
}

socket.onmessage = function(event) {
    const data = JSON.parse(event.data)
    console.log(data);
}

socket.onclose = function(event) {
    if (event.wasClean) {
        console.log("connection closed.");
    } else {
        console.log("connection died.");
    }
}

socket.onerror = function(error) {
    console.log(error);
}
