class CreateDomains < ActiveRecord::Migration[5.0]
  def change
    create_table :domains do |t|
      t.string :name
      t.string :type
      t.string :ttl
      t.string :notes
      t.string :primary_ns
      t.string :contact
      t.integer :refresh
      t.integer :retry
      t.integer :expire
      t.integer :minimum
      t.string :authority_type, limit: 1

      t.timestamps
    end
  end
end
