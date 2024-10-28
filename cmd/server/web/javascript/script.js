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
