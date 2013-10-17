class Login < ActiveRecord::Base
  belongs_to :app
  belongs_to :identity

  validates :email,    presence: true

  before_create :generate_authcode
  before_create :set_identity
  after_create  :deliver_authcode_email

  state_machine :status, :initial => :requested do
    event :mark_confirmed! do
      transition any => :confirmed
    end
  end

  def self.confirm(params={})
    app   = App.where(app_name: params[:app_name]).first
    return [nil, "Sorry, we couldn't find an app by that app_name."] if !app

    login = app.logins.where(email: params[:email]).first
    return [nil, "Sorry, we couldn't find a login request using that email."] if !login
    return [nil, "Sorry, the authcode is incorrect."] if login.authcode != params[:authcode]

    if login.requested? && login.mark_confirmed!
      [login]
    else
      [nil, "Sorry, this authcode has already been used with this email."]
    end
  end

  private

  def generate_authcode
    self.authcode = random_authcode 
  end

  def random_authcode
    #'%04d' % rand(10 ** 4)
    rand.to_s[2..5]
  end

  def set_identity
    identity = app.identities.where(email: email).first

    if identity
      self.identity = identity
    else
      self.identity = app.identities.create(email: email)
    end
  end

  def deliver_authcode_email 
    puts "="*100
    response = Pony.mail({
      to:       email,
      from:     FROM,
      subject:  "Your code: #{authcode}. Please enter it to login.",
      body:     "Your code: #{authcode}. Please enter it to login.",
      via:      :smtp,
      :via_options => {
        :address        => SMTP_ADDRESS,
        :port           => SMTP_PORT,
        :user_name      => SMTP_USERNAME,
        :password       => SMTP_PASSWORD,
        :authentication => :plain
      }
    })

    puts response
    puts "="*100
  end
end
