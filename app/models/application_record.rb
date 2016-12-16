class ApplicationRecord < ActiveRecord::Base
  self.abstract_class = true
  self.inheritance_column = '__unknown__' # don't use "type" as inheritance column
end
