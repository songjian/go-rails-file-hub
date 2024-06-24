class AddOnlineToClients < ActiveRecord::Migration[7.1]
  def change
    add_column :clients, :online, :boolean, default: false
  end
end
