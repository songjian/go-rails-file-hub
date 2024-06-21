ActiveAdmin.register Client do
  permit_params :name, :api_key, :active, :description
  
  index do
    selectable_column
    id_column
    column :name
    column :api_key
    column :active
    column :description
    actions
  end

  filter :name
  filter :api_key
  filter :active
  filter :description

  form do |f|
    f.inputs do
      f.input :name
      f.input :api_key
      f.input :active
      f.input :description
    end
    f.actions
  end
  
end
