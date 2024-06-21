class CreateRemoteFiles < ActiveRecord::Migration[7.1]
  def change
    create_table :remote_files do |t|
      t.string :filename
      t.string :path
      t.string :file_type
      t.datetime :file_created_at
      t.integer :file_size
      t.datetime :last_modified
      t.references :client, null: false, foreign_key: true

      t.timestamps
    end
  end
end
