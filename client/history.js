import { getGraphEntities, processAndAppendServerData } from "./common.js";

const interfaceOrAdapterName = localStorage.getItem('iface');

// The canvas
const graphDiv = document.getElementById('graphDiv');

// Refer: https://developers.google.com/chart/interactive/docs/gallery/linechart
const rawData = {
    received: [0],
    sent: [0],
    prev: {}
}

fetch(`http://${window.location.host}/getHistory?iface=${interfaceOrAdapterName}`)
.then((response) => response.json())
.then((data) => {

    google.charts.load('current', {packages: ['corechart', 'line']});
    google.charts.setOnLoadCallback(() => {


        const graphEntites = getGraphEntities();
        const graph = graphEntites[0];
        const graphData = graphEntites[1];
        const graphOptions = graphEntites[2];

        for (const dataPoint of data) {
            processAndAppendServerData(rawData, dataPoint)
            const index = rawData.received.length - 1;
            const received = rawData.received[index];
            const sent = rawData.sent[index];
            graphData.addRows([[index, sent, received]]);
        }
        graph.draw(graphData, graphOptions);
    });
})
.catch((err) => {
    console.log(err);
})
