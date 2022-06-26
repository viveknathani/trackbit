export function processAndAppendServerData(rawData, serverData) {

    if (Object.keys(rawData.prev).length !== 0) {
        const receivedKiloBytes = (serverData.BytesReceived - rawData.prev.BytesReceived) / 1024;
        const sentKiloBytes = (serverData.BytesSent - rawData.prev.BytesSent) / 1024;
        rawData.sent.push(receivedKiloBytes);
        rawData.received.push(sentKiloBytes);
    }
    rawData.prev = serverData;
}

export function getGraphEntities() {

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
    const graphDiv = document.getElementById('graphDiv');
    const graph = new google.visualization.LineChart(graphDiv);
    graphData.addColumn('number', 'Seconds');
    graphData.addColumn('number', 'Received');
    graphData.addColumn('number', 'Sent');
    graphData.addRows([[0, 0, 0]]);

    return [graph, graphData, graphOptions]
}
    