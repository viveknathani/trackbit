import { getGraphEntities, processAndAppendServerData } from "./common.js";

const MAX_VISIBLE_SECONDS = 60;
const interfaceOrAdapterName = prompt("Enter interface or adapter name: ")

// Store data coming in from the server
const rawData = {
    received: [0],
    sent: [0],
    prev: {}
}

// Refer: https://developers.google.com/chart/interactive/docs/gallery/linechart
google.charts.load('current', {packages: ['corechart', 'line']});
google.charts.setOnLoadCallback(run);

function run() {

    const graphEntites = getGraphEntities();
    const graph = graphEntites[0];
    const graphData = graphEntites[1];
    const graphOptions = graphEntites[2];
    const socket = new WebSocket(`ws://${window.location.host}/stats`);
    socket.onopen = onSocketOpen;
    socket.onerror = onSocketError;
    socket.onclose = onSocketClose;
    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        processAndAppendServerData(rawData, data);
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

// Utility functions for websockets

function onSocketOpen(event) {
    console.log('connection established.');

    // Here, "this" is the WebSocket instance.
    localStorage.setItem("iface", interfaceOrAdapterName)
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

