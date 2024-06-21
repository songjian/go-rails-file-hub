class RemoteFile < ApplicationRecord
  belongs_to :client

  def self.ransackable_attributes(auth_object = nil)
    ["client_id", "created_at", "file_created_at", "file_size", "file_type", "filename", "id", "id_value", "last_modified", "path", "updated_at"]
  end
end
