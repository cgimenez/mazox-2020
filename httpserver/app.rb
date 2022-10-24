require "bundler"
Bundler.require

require "sinatra/base"
# require "sinatra/custom_logger"

disable :logging
use Rack::CommonLogger
use Rack::Deflater

set :port, 3010
set :bind, "0.0.0.0"
set :sockets, []
set :root, __dir__ + "/../www"

set :public_folder, __dir__ + "/../www"

get "/" do
  send_file "../www/index.html"
end

get "/ws" do
  request.websocket do |ws|
    ws.onopen do
      settings.sockets << ws
    end

    ws.onmessage do |msg|
      script_file = File.dirname(__FILE__) + "/../www/mazox.wasm"
      ft = File.mtime(script_file)
      timer = EventMachine::PeriodicTimer.new(0.5) do
        if File.mtime(script_file) != ft
          ft = File.mtime(script_file)
          logger.info "Send reload"
          settings.sockets.each { |s| s.send("reload") }
          timer.cancel
        end
      end
    end

    ws.onclose do
      logger.info "Websocket closed"
      settings.sockets.delete(ws)
    end
  end
end
