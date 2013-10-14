class Identity < ActiveRecord::Base
  belongs_to  :app
  has_many    :logins

  validates   :email, presence: true, uniqueness: true

  before_save :set_short_id

  private

  def set_short_id
    self.short_id = short_id || ["IDNT_", Shorty.next_short_id(Identity)].join
  end
end
