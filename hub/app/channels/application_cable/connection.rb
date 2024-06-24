module ApplicationCable
  class Connection < ActionCable::Connection::Base
    identified_by :current_client


    def connect
      self.current_client = find_verified_client
      current_client.update(online: true) if current_client
    end

    def disconnect
      current_client.update(online: false) if current_client
    end

    private
      def find_verified_client
        bearer_token = request.headers['Authorization']
        if bearer_token.present? && bearer_token.start_with?('Bearer ')
          api_key = bearer_token.split(' ').last
          if verified_client = Client.find_by(api_key: api_key)
            verified_client
          else
            reject_unauthorized_connection
          end
        else
          reject_unauthorized_connection
        end
      end

  end
end
