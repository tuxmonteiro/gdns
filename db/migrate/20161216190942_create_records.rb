class CreateRecords < ActiveRecord::Migration[5.0]
  def change
    create_table :records do |t|
      t.string :name
      t.string :type
      t.string :content
      t.references :domain, foreign_key: true

      t.timestamps
    end
  end
end
