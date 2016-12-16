Rails.application.routes.draw do
  post 'bind9/schedule_export'
  post 'bind9/export'

  resources :domains do
    resources :records
  end

end
