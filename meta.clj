#!/usr/bin/env bb

(ns script)

(defn archive
  "Create an archive, combining all the sources into the destination. The type
  of archive is determined by the destination extension."
  [sources destination opts])

(defn meta-str
  [sym]
  (-> (meta sym)
      (update :name str)
      (select-keys [:name :doc :arglists])
      pr-str))


#_(meta-str #'archive)
