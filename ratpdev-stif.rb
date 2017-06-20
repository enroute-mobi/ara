# coding: utf-8
require 'rest-client'
require 'json'

server = 'http://stif1.api.concerto.ratpdev.com'
server = 'http://localhost:8080'


begin
  referentials = JSON.parse(RestClient.get "#{server}/_referentials")

  referential = referentials.find { |r| r["Slug"] == "concerto" }
  if referential
    puts "✓ Remove concerto referential"
    RestClient.delete "#{server}/_referentials/#{referential['Id']}"
  end

  attributes = { "Slug" => "concerto" }
  puts "✓ Create concerto referential"
  RestClient.post "#{server}/_referentials", attributes.to_json, {content_type: :json}

  ineo_partner = {
    "slug" => "sqybus",
    "connectorTypes" => %w{siri-stop-monitoring-request-collector siri-check-status-client},
    "settings" => {
      remote_credential: "RATPDEV:Concerto",
      remote_objectid_kind: "hastus",
      remote_url: "http://194.206.198.37:8085/ProfilSiriKidf2_4Producer-Sqybus/SiriServices"
    }
  }

  puts "✓ Create ineo partner"
        RestClient.post "#{server}/concerto/partners", ineo_partner.to_json, {content_type: :json}

  stif_partner = {
    "slug" => "stif",
    "connectorTypes" => %w{siri-check-status-server siri-check-status-client},
    "settings" => {
      local_credential: "RELAIS",
      remote_objectid_kind: "stif",
      remote_credential: "RATPDev",
      remote_url: "http://emission.recette.relais-ivtr.stif.info:8080/emission/SiriRelais/SiriProducerRpcPort"
    }
  }

  puts "✓ Create stif partner"
        RestClient.post "#{server}/concerto/partners", stif_partner.to_json, {content_type: :json}

  concerto_stop_areas = [
     ['booarle', 'STIF:StopPoint:Q:1:'],
     ['boabonn', 'STIF:StopPoint:Q:2:'],
     ['boaclair', 'STIF:StopPoint:Q:3:'],
     ['boacroi', 'STIF:StopPoint:Q:4:'],
     ['boaegli', 'STIF:StopPoint:Q:5:'],
     ['boahote', 'STIF:StopPoint:Q:6:'],
     ['boalang', 'STIF:StopPoint:Q:7:'],
     ['boapvc', 'STIF:StopPoint:Q:8:'],
     ['boatati', 'STIF:StopPoint:Q:9:'],
     ['boavons', 'STIF:StopPoint:Q:10:']
  ]

  concerto_stop_areas.each do |values|
    attributes = {
      "ObjectIDs" => { "hastus" => values[0] }
      # "ObjectIDs" => { "hastus" => values[0], "stif" => values[1] }
    }
    puts "✓ Create Concerto StopArea '#{values[0]}'"
    RestClient.post "#{server}/concerto/stop_areas", attributes.to_json, {content_type: :json}
  end
rescue RestClient::ExceptionWithResponse => err
  puts err.response.body
  raise err
end
