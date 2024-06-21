class FileChannel < ApplicationCable::Channel
  def subscribed
    stream_from "file_channel_#{current_client.id}"
  end

  def receive(data)
    puts "Data: #{data.inspect}"
    if data['file_action'] == 'created'
      ActiveRecord::Base.transaction do
        RemoteFile.create(filename: data['filename'], path: data['path'], file_type: data['file_type'], file_created_at: data['operation_time'], file_size: data['file_size'], last_modified: data['operation_time'], client_id: current_client.id)
      end
    elsif data['file_action'] == 'deleted'
      RemoteFile.find_by(path: data['path'], filename: data['filename'], client_id: current_client.id).destroy
    end
  end
  
end
