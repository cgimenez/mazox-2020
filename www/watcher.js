window.addEventListener("load", (e) => {
  function reload() {
    console.log('reload !')
    window.location.reload(true)
  }

  // const socket = new WebSocket('ws://127.0.0.1:3010/ws')
  const socket = new WebSocket("ws://" + location.host + "/ws")

  socket.addEventListener('open', function (event) {
    socket.send('voidz')
  })

  socket.addEventListener('close', function (event) {
    //setTimeout(function () { reload(); }, 1000);
    reload()
  })

  socket.addEventListener('message', function (event) {
    if (event.data === 'reload') {
      reload()
    }
  })
})