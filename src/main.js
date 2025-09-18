import ReconnectingWebSocket from "reconnecting-websocket";

let secondsSinceLastPress = 0;
let ws = null;
let isConnected = false;

function initWebSocket() {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsUrl = `${protocol}//${window.location.host}/ws`;

  ws = new ReconnectingWebSocket(wsUrl);

  ws.addEventListener("open", () => {
    console.log("WebSocket connected");
    isConnected = true;
    updateDisplay();
  });

  ws.addEventListener("message", (event) => {
    const message = JSON.parse(event.data);
    if (message.type === "update") {
      secondsSinceLastPress = message.counter;
      updateDisplay();
    }
  });

  ws.addEventListener("close", () => {
    console.log("WebSocket disconnected");
    isConnected = false;
    updateDisplay();
  });

  ws.addEventListener("error", (error) => {
    console.error("WebSocket error:", error);
  });
}

function updateDisplay() {
  const button = document.getElementById("counter-button");
  if (button) {
    if (isConnected) {
      button.textContent = `${secondsSinceLastPress} seconds since the last press`;
      button.classList.remove("disconnected");
    } else {
      button.textContent = "disconnected";
      button.classList.add("disconnected");
    }
  }
}

function resetCounter() {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: "reset" }));
  }
}

document.addEventListener("DOMContentLoaded", () => {
  initWebSocket();

  const button = document.getElementById("counter-button");
  if (button) {
    button.addEventListener("click", resetCounter);
    updateDisplay();
  }
});
