import EventEmitter from 'events';

// Takes a WebSocket or RTCDataChannel and implements a simple RPC protocol
class RpcSocket extends EventEmitter {
  constructor(socket) {
    super();
    this._m = Object.create(null);
    this._i = 0;
    this._s = socket;
    if (socket.readyState === socket.OPEN) {
      this._b = null;
    } else {
      this._b = [];
      this._s.onopen = () => {
        this._b.forEach(data => {
          this._s.send(data);
        });
        this._b = null;
        this.emit('open');
      };
    }
    this._s.onclose = () => {
      Object.keys(this._m).forEach((key) => {
        this._m[key][1].call(this);
        delete this._m[key];
      });
      this.emit('close');
    };
    this._s.onerror = () => {
      this.emit('error');
    };
    this._s.onmessage = (e) => {
      const obj = JSON.parse(e.data);
      if ('method' in obj) {
        if ('id' in obj) {
          const resolve = (result) => {
            this._s.send(JSON.stringif({
              id: obj.id,
              result: result,
            }));
          };
          const reject = (error) => {
            this._s.send(JSON.stringif({
              id: obj.id,
              error: error,
            }));
          };
          this.emit('call', resolve, reject, obj.params);
        } else {
          this.emit('notify', obj.params);
        }
      } else if ('result' in obj) {
        if (typeof obj.id !== 'number') {
          return;
        }
        this._m[obj.id][0].call(this, obj.result);
        delete this._m[obj.id];
      } else if ('error' in obj) {
        if (typeof obj.id !== 'number') {
          return;
        }
        this._m[obj.id][1].call(this, obj.error);
        delete this._m[obj.id];
      }
    };
  }

  _send(data) {
    if (this._b !== null) {
      this._b.push(data);
    } else {
      this._s.send(data);
    }
  }

  call(method, params) {
    this._i += 1;
    const id = this._i;
    this._send(JSON.stringify({
      id: id,
      method: method,
      params: params,
    }));
    return new Promise((resolve, reject) => {
      this._m[id] = [resolve, reject];
    });
  }

  notify(method, params) {
    this._send(JSON.stringify({
      method: method,
      params: params,
    }));
  }
}

export default RpcSocket;
