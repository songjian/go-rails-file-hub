class CreateClients < ActiveRecord::Migration[7.1]
  def change
    create_table :clients do |t|
      t.string :name
      t.string :api_key
      t.boolean :active, default: true
      t.text :description

      t.timestamps
    end
  end
end
