class Client < ApplicationRecord
  after_initialize :set_defaults, if: :new_record?

  def self.ransackable_attributes(auth_object = nil)
    ["active", "api_key", "created_at", "description", "id", "id_value", "name", "updated_at"]
  end

  private

  def set_defaults
    self.api_key ||= SecureRandom.hex
  end
end
