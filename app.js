window.onload = function() {
  
    // get the references of the page elements.
    var form = document.getElementById('form-msg');
    var txtMsg = document.getElementById('msg');
    var listMsgs = document.getElementById('msgs');
    var socketStatus = document.getElementById('status');
    var btnClose = document.getElementById('close');
  };

  // Creating a new WebSocket connection.
 var socket = new WebSocket('ws://echo.websocket.org');

 socket.onopen = function(event) {
	socketStatus.innerHTML = 'Connected to: ' + event.currentTarget.URL;
	socketStatus.className = 'open';
};

socket.onerror = function(error) {
	console.log('WebSocket error: ' + error);
};

socket.send(data); 

form.onsubmit = function(e) {
    e.preventDefault();
    
    // Recovering the message of the textarea.
    var msg = txtMsg.value;
    
    // Sending the msg via WebSocket.
    socket.send(msg);
    
    // Adding the msg in a list of sent messages.
    litsMsgs.innerHTML += '<li class="sent"><span>Sent:</span>' + msg + '</li>';
    
    // Cleaning up the field after sending.
    txtMsg.value = '';
    
    return false;
};

socket.onmessage = function(event) {
    var msg = event.data;
    listMsgs.innerHTML += '<li class="received"><span>Received:</span>' + msg + '</li>';
};

socket.close();

socket.onclose = function(event) {
    socketStatus.innerHTML = 'Disconnected from the WebSocket.';
    socketStatus.className = 'closed';
};