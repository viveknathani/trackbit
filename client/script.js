const MAX_VISIBLE_SECONDS = 60;
const interfaceOrAdapterName = prompt("Enter interface or adapter name: ")

// Store data coming in from the server
const rawData = {
    received: [0],
    sent: [0],
    prev: {}
}

// The canvas
const graphDiv = document.getElementById('graphDiv');

// Refer: https://developers.google.com/chart/interactive/docs/gallery/linechart
google.charts.load('current', {packages: ['corechart', 'line']});
google.charts.setOnLoadCallback(run);

function run() {

    const graphData = new google.visualization.DataTable();
    const graphOptions = {
        width: 800,
        height: 300,
        hAxis: {
            title: 'Seconds'
        },
        vAxis: {
            title: 'KB/s',
            minValue: 0
        }
    };
    const graph = new google.visualization.LineChart(graphDiv);
    graphData.addColumn('number', 'Seconds');
    graphData.addColumn('number', 'Received');
    graphData.addColumn('number', 'Sent');
    graphData.addRows([[0, 0, 0]]);


    const socket = new WebSocket(`ws://${window.location.host}/stats`);
    socket.onopen = onSocketOpen;
    socket.onerror = onSocketError;
    socket.onclose = onSocketClose;
    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        processAndAppendServerData(data);
        const index = rawData.received.length - 1;
        const received = rawData.received[index];
        const sent = rawData.sent[index];
        graphData.addRows([[index, sent, received]]);
        if (graphData.getNumberOfRows() >= MAX_VISIBLE_SECONDS) {
            graphData.removeRow(0);
        }
        graph.draw(graphData, graphOptions);
    }
}

// Server sends data via a WebSocket connection.
function processAndAppendServerData(serverData) {

    if (Object.keys(rawData.prev).length !== 0) {
        const receivedKiloBytes = (serverData.BytesReceived - rawData.prev.BytesReceived) / 1024;
        const sentKiloBytes = (serverData.BytesSent - rawData.prev.BytesSent) / 1024;
        rawData.sent.push(receivedKiloBytes);
        rawData.received.push(sentKiloBytes);
    }
    rawData.prev = serverData;
}

// Utility functions for websockets

function onSocketOpen(event) {
    console.log('connection established.');

    // Here, "this" is the WebSocket instance.
    this.send(JSON.stringify({
        Name: interfaceOrAdapterName
    }));
}

function onSocketClose(event) {
    if (event.wasClean) {
        console.log("connection closed.");
    } else {
        console.log("connection died.");
    }
}

function onSocketError(error) {
    console.log(error);
}

