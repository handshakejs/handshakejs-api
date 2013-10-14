class Shorty
  def self.next_short_id(klass, hex=20)
    safety = 0
    begin
      safety += 1
      result = SecureRandom.hex(hex).upcase
    end while klass.exists?(short_id: result) && safety < 5000 # only attempt 5000 times
    result
  end
end
