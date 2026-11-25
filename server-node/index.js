const express = require('express');
const http = require('http');
const WebSocket = require('ws');

const app = express();
app.use(express.static('../web'));

const server = http.createServer(app);
const wss = new WebSocket.Server({ server, path: '/ws' });

const rooms = {};
function getRoom(id, owner) {
  if (!rooms[id]) rooms[id] = { id, ownerId: owner, thLevel: 1, players: {} };
  return rooms[id];
}

wss.on('connection', (ws) => {
  ws.on('message', (raw) => {
    try {
      const cmd = JSON.parse(raw);
      ws.send(JSON.stringify({ received: cmd }));
    } catch (e) {
      ws.send(JSON.stringify({ error: 'bad json' }));
    }
  });
});

app.get('/api/room/create', (req, res) => {
  const { id, owner } = req.query;
  if (!id || !owner) return res.status(400).send('missing id/owner');
  const base = getRoom(id, owner);
  if (!base.players[owner]) base.players[owner] = { resources: 1000, builders: 2, assigned: true };
  res.json(base);
});

app.get('/api/room/join', (req, res) => {
  const { id, player } = req.query;
  const base = rooms[id];
  if (!base) return res.status(404).send('room not found');
  if (Object.keys(base.players).length >= 4) return res.status(400).send('room full');
  if (!base.players[player]) base.players[player] = { resources: 1000, builders: 2, assigned: false };
  res.json(base);
});

app.get('/api/room/upgrade_th', (req, res) => {
  const { id } = req.query;
  const base = rooms[id];
  if (!base) return res.status(404).send('room not found');
  if (base.thLevel < 16) base.thLevel++;
  res.json(base);
});

const port = process.env.PORT || 8080;
server.listen(port, () => console.log('server listening on :' + port));
