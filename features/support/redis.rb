require "redis"

Before do
  if ENV["ARA_REDIS_ADDR"]
    ENV["ARA_REDIS_CODESPACES"] = "internal,external" unless ENV["ARA_REDIS_CODESPACES"]
  end
end

After do
  if ENV["ARA_REDIS_ADDR"]
    redis = Redis.new
    begin
        redis.ping
        redis.flushall
    rescue Redis::BaseError => e
        e.inspect
        e.message
    end
  end
end