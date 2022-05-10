require 'siri/xsd'

# Helper to validate a XML content and log them with details
module Siri
  class Validator

    def self.enabled?
      @enabled ||= (ENV['SIRI_VALIDATE'] == "true")
    end

    def self.create(content, description = nil)
      if enabled?
        new content, description
      else
        Null.instance
      end
    end

    class Null
      include Singleton

      def errors; []; end
      def log; end
    end

    def initialize(content, description = nil)
      @content = content
      @description = description
    end
    attr_reader :content, :description

    def validator
      @validator ||= Siri::Xsd::Validator.new
    end

    def validate
      unless validated?
        validator.validate(StringIO.new(content))
        @validated = true
      end
    end

    def validated?
      @validated ||= false
    end

    def errors
      validate
      validator.errors
    end

    def in_description
      " in #{description}" if description
    end

    def log
      unless errors.empty?
        logger.log "✗ #{errors.count} XSD error(s)#{in_description}"
        logger.log content
        errors.each do |error|
          logger.log "#{error.to_s}"
        end
      else
        logger.log "✓ No XSD error#{in_description}"
      end
    end

    class << self
      attr_accessor :logger
    end

    def logger
      self.class.logger
    end
  end
end

Before do
  Siri::Validator.logger = self
end
