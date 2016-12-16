class RecordSerializer < ActiveModel::Serializer
  type 'record'
  attributes :id, :name, :type, :content
  has_one :domain
end
