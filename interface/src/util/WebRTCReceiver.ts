
export class WebRTCReceiver {
  private terminated: boolean;
  private ws: WebSocket | null;
  private pc: RTCPeerConnection | null;
  private restartTimeout: number | null;
  private restartPause: number = 2000;
  private name: string;

  constructor(url: string, name: string) {
    this.terminated = false;
    this.name = name;
    this.ws = null;
    this.pc = null;
    this.restartTimeout = null;
    this.start(`${url}ws`);
  }

  start(url: string) {
    console.log('connecting');

    const fullUrl = url;

    this.ws = new WebSocket(fullUrl);

    this.ws.onerror = (e: Event) => {
      console.log('ws error', e);
      if (this.ws === null) {
        return;
      }
      this.ws.close();
      this.ws = null;
    };

    this.ws.onclose = () => {
      console.log('ws closed');
      this.ws = null;
      this.scheduleRestart(fullUrl);
    };

    this.ws.onmessage = (msg) => this.onIceServers(msg, fullUrl);
  }

  onIceServers(msg: MessageEvent, url: string) {
    if (this.ws === null) {
      return;
    }

    const iceServers = JSON.parse(msg.data);

    this.pc = new RTCPeerConnection({
      iceServers,
    });

    this.ws.onmessage = (msg) => this.onRemoteDescription(msg);
    this.pc.onicecandidate = (evt) => this.onIceCandidate(evt);

    this.pc.oniceconnectionstatechange = () => {
      if (this.pc === null) {
        return;
      }

      console.log('peer connection state:', this.pc.iceConnectionState);

      switch (this.pc.iceConnectionState) {
        case 'disconnected':
          this.scheduleRestart(url);
      }
    };

    this.pc.ontrack = (evt) => {
      console.log('new track ' + evt.track.kind);
      const video = document.getElementById(this.name) as HTMLVideoElement;
      if (!video) return;
      video.srcObject = evt.streams[0];

      video.addEventListener('loadedmetadata', () => {
        console.log(`Stream dimensions: ${video.videoWidth}x${video.videoHeight}`);
      })
    };

    const direction = 'sendrecv';
    this.pc.addTransceiver('video', { direction });
    this.pc.addTransceiver('audio', { direction });

    this.pc.createOffer().then((desc) => {
      if (this.pc === null || this.ws === null) {
        return;
      }

      this.pc.setLocalDescription(desc);

      console.log('sending offer');
      this.ws.send(JSON.stringify(desc));
    });
  }

  onRemoteDescription(msg: MessageEvent) {
    if (this.pc === null || this.ws === null) {
      return;
    }

    this.pc.setRemoteDescription(
      new RTCSessionDescription(JSON.parse(msg.data))
    );
    this.ws.onmessage = (msg) => this.onRemoteCandidate(msg);
  }

  onIceCandidate(evt: RTCPeerConnectionIceEvent) {
    if (this.ws === null) {
      return;
    }

    if (evt.candidate !== null) {
      if (evt.candidate.candidate !== '') {
        this.ws.send(JSON.stringify(evt.candidate));
      }
    }
  }

  onRemoteCandidate(msg: MessageEvent) {
    if (this.pc === null) {
      return;
    }

    this.pc.addIceCandidate(JSON.parse(msg.data));
  }

  scheduleRestart(url: string) {
    if (this.terminated) {
      return;
    }

    if (this.ws !== null) {
      this.ws.close();
      this.ws = null;
    }

    if (this.pc !== null) {
      this.pc.close();
      this.pc = null;
    }

    this.restartTimeout = window.setTimeout(() => {
      this.restartTimeout = null;
      this.start(url);
    }, this.restartPause);
  }
}
