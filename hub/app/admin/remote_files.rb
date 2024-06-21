ActiveAdmin.register RemoteFile do
  permit_params :filename, :path, :file_type, :file_created_at, :file_size, :last_modified, :client_id

  index do
    selectable_column
    id_column
    column :filename
    column :path
    column :file_type
    column :file_size
    column :file_created_at
    column :last_modified
    column :client
    actions
  end

  filter :filename
  filter :path
  filter :file_type
  filter :file_size
  filter :file_created_at
  filter :last_modified

  form do |f|
    f.inputs do
      f.input :filename
      f.input :path
      f.input :file_type
      f.input :file_created_at
      f.input :client
    end
    f.actions
  end
end
