const clientListContainer = document.getElementById('clientList');

function fetchClients() {
    fetch('/aliveConnection')
    .then(response => response.json())
    .then(clients => {
        clientListContainer.innerHTML = '';
        clients.forEach(client => {
            const clientItem = document.createElement('div');
            clientItem.className = 'list-group-item d-flex justify-content-between align-items-center';
            clientItem.innerHTML = `
            <span>${client.ip}</span>
            <button class="btn btn-danger btn-sm" onclick="shutdownClient('${client.ip}')">Spegni</button>
            `;
            clientListContainer.appendChild(clientItem);
        });
    })
    .catch(error => {
        clientListContainer.innerHTML = '<div class="text-danger">Failed to load clients.</div>';
    });
}

function shutdownClient(ip) {
    fetch('/shutdownClient', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ip })
    })
    .then(response => response.json())
    .then(data => {
        alert(`Client ${ip} shutdown: ${data.message}`);
        fetchClients(); // Refresh client list after shutdown
    })
    .catch(error => {
        alert(`Failed to shutdown client ${ip}`);
    });
}

// Update the client list every 3 seconds
setInterval(fetchClients, 3000);
fetchClients();

document.getElementById("powerOnButton").onclick = function() {
    fetch('/powerOn', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById("resultPowerOn").innerText = data.message;
        document.getElementById("resultPowerOn").classList.add('text-success');
        document.getElementById("resultPowerOff").innerText = '';
    })
    .catch(error => {
        document.getElementById("resultPowerOn").innerText = 'Failed to power on!';
        document.getElementById("resultPowerOn").classList.add('text-danger');
    });
};

document.getElementById("powerOffButton").onclick = function() {
    fetch('/powerOff', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById("resultPowerOff").innerText = data.message;
        document.getElementById("resultPowerOff").classList.add('text-danger');
        document.getElementById("resultPowerOn").innerText = '';
    })
    .catch(error => {
        document.getElementById("resultPowerOff").innerText = 'Failed to power off!';
        document.getElementById("resultPowerOff").classList.add('text-danger');
    });
};
