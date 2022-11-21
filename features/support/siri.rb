require 'siri/xsd'

# Helper to validate a XML content and log them with details
module Siri
  class Validator

    def self.enabled?
      @enabled ||= %w{true strict}.include?(ENV['SIRI_VALIDATE'])
    end

    mattr_accessor :error_count, default: 0

    def self.strict_mode?
      @strict_mode ||= (ENV['SIRI_VALIDATE'] == "strict")
    end

    def self.enable
      previous_value = enabled?

      @enabled = true
      yield if block_given?
    ensure
      @enabled = previous_value
    end

    def self.detect_error?(&block)
      before_error_count = self.error_count
      enable(&block)
      before_error_count != self.error_count
    end

    def self.fail_on_error(scenario, block)
      if detect_error? { block.call }
        scenario.fail("Invalid SIRI XML detected") 
      end
    end

    def self.create(content, description = nil)
      if enabled?
        new content, description: description
      else
        Null.instance
      end
    end

    class Null
      include Singleton

      def errors; []; end
      def log; end
    end

    def initialize(content, description: nil)
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

        self.error_count += validator.errors.count
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

Around do |scenario, block|
  unless Siri::Validator.strict_mode?
    block.call
  else
    Siri::Validator.fail_on_error(scenario, block)
  end
end

Around('@siri-valid') do |scenario, block|
  Siri::Validator.fail_on_error(scenario, block)
end