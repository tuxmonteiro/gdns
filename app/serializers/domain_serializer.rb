class DomainSerializer < ActiveModel::Serializer
  type 'domain'
  attributes :id, :name, :type, :ttl, :notes, :primary_ns, :contact, :refresh, :retry, :expire, :minimum, :authority_type
end
