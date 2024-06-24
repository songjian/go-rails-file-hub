class RemoteFile < ApplicationRecord
  belongs_to :client

  def self.ransackable_attributes(auth_object = nil)
    ["client_id", "created_at", "file_created_at", "file_size", "file_type", "filename", "id", "id_value", "last_modified", "path", "updated_at"]
  end

  # Delete remote file
  def delete_file
    data = {file_action: 'delete', filename: self.filename, path: self.path, file_type: self.file_type}
    ActionCable.server.broadcast("file_channel_#{self.client_id}", data)
  end

end
