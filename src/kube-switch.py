#!/usr/bin/env python3
#
# The MIT License (MIT)
# Copyright 2020, Rafael Grigorian <rafael@grigorian.org>
#
# ---
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import yaml, re, sys, tty, termios
from os import path
from kubernetes import client, config

def get_char ():
	fd = sys.stdin.fileno ()
	old_settings = termios.tcgetattr ( fd )
	try:
		tty.setcbreak ( fd )
		ch = sys.stdin.read ( 1 )
	finally:
		termios.tcsetattr ( fd, termios.TCSADRAIN, old_settings )
	return ch

class ANSI:

	def civis ():
		print ( "\033[?25l", end = "" )

	def cnorm ():
		print ( "\033[34h\033[?25h", end = "" )

	def smcup ():
		print ( "\033[?1049h\033[H", end = "" )

	def rmcup ():
		print ( "\033[?1049l", end = "" )

	def clear ():
		print ( "\033[H\033[2J", end = "" )

	def set_cursor ( r, c ):
		print ( "\033[%d;%dH" % ( r, c ), end = "" )

	def selection_color ( message ):
		return "\033[43;30m%s\033[0m" % message

	def active_color ( message ):
		return "\033[33m%s\033[0m" % message

	def error_color ( message ):
		return "\033[31m%s\033[0m" % message

def is_arrow ( queue ):
	if len ( queue ) == 3:
		return queue [ 0 ] == 27 and queue [ 1 ] == 91
	return False

class Kube:

	config_path = None
	config_yaml = None

	contexts = None
	context = None
	context_names = None
	active_context = None
	active_context_index = None

	namespaces = None
	namespace = None
	namespace_names = None
	active_namespace = None
	active_namespace_index = None

	def __init__ ( self ):
		relative = config.kube_config.KUBE_CONFIG_DEFAULT_LOCATION
		self.config_path = path.expanduser ( relative )
		self.load_config ()
		self.contexts, self.context = config.list_kube_config_contexts ()
		self.context_names = list ( map ( lambda x : x ["name"], self.contexts ) )
		self.active_context = self.context ["name"]
		self.active_namespace = self.context ["context"] ["namespace"]
		self.active_context_index = self.context_names.index ( self.active_context )

	def load_config ( self ):
		with open ( self.config_path ) as file:
			self.config_yaml = yaml.load ( file, Loader = yaml.FullLoader )

	def save_config ( self ):
		with open ( self.config_path, 'w' ) as file:
			yaml.dump ( self.config_yaml, file, default_flow_style = False )

	def get_api ( self, context_name ):
		api_config = config.new_client_from_config ( context = context_name )
		return client.CoreV1Api ( api_client = api_config )

	def set_context ( self, chosen ):
		self.config_yaml ["current-context"] = chosen
		self.save_config ()
		self.context = next ( e for e in self.contexts if e ["name"] == chosen )
		self.active_context = self.context ["name"]
		self.active_namespace = self.context ["context"] ["namespace"]
		self.active_context_index = self.context_names.index ( self.active_context )
		self.namespaces = list ( self.get_api ( chosen ).list_namespace ( timeout_seconds = 1 ).items )
		self.namespace = next ( e for e in self.namespaces if e.metadata.name == self.active_namespace )
		self.namespace_names = list ( map ( lambda x : x .metadata.name, self.namespaces ) )
		self.active_namespace_index = self.namespace_names.index ( self.active_namespace )

	def set_namespace ( self, chosen ):
		self.namespace = next ( e for e in self.namespaces if e.metadata.name == chosen )
		self.active_namespace = chosen
		self.active_namespace_index = self.namespace_names.index ( chosen )
		index = list ( map ( lambda x : x ["name"], self.contexts ) ).index ( self.active_context )
		self.config_yaml ["contexts"] [ index ] ["context"] ["namespace"] = chosen
		self.save_config ()

	def prompt ( self, title, choices, active_index = None ):
		filter = ""
		key_queue = []
		selection = active_index
		current_value = choices [ active_index ]
		return_value = None
		while True:
			filtered_choices = [i for i in choices if filter in i]
			selection = filtered_choices.index ( current_value ) if selection == -2 and current_value in filtered_choices else selection
			selection = min ( max ( selection, 0 ), len ( filtered_choices ) - 1 )
			selected_value = None if len ( filtered_choices ) == 0 else filtered_choices [ selection ]
			ANSI.clear ()
			ANSI.civis ()
			ANSI.set_cursor ( 1, 1 )
			print ( "Filter: %s" % filter, end = "" )
			print ( ANSI.active_color ("_") )
			print ("Type to filter, UP/DOWN move, ENTER select, ESC exit\n")
			print ( "%s: %s\n" % ( title, selected_value ) )
			for index, choice in enumerate ( filtered_choices ):
				if index == selection:
					print ( ANSI.selection_color ( " + %s " % choice ) )
				elif current_value == choice:
					print ( ANSI.active_color ( " + %s " % choice ) )
				else:
					print ( " + %s" % choice )
			if len ( filtered_choices ) == 0:
				print ( ANSI.error_color ("No results; refine filter.") )
			ANSI.set_cursor ( 1, len ( filter ) + 9 )
			char = get_char ()
			key_queue.append ( ord ( char ) )
			key_queue = key_queue [-3:]
			ANSI.set_cursor ( 30, 1 )
			if ord ( char ) == 3:
				return_value = None
				break
			elif ord ( char ) == 13 or ord ( char ) == 10:
				if len ( filtered_choices ) != 0:
					return_value = selected_value
					break
			elif ord ( char ) == 127:
				filter = filter [0:-1]
				selection = -2
			elif ord ( char ) == 68 or ord ( char ) == 67:
				continue
			elif is_arrow ( key_queue ) and ord ( char ) == 65:
				if len ( filtered_choices ) > 0:
					selection -= 1
			elif is_arrow ( key_queue ) and ord ( char ) == 66:
				if len ( filtered_choices ) > 0:
					selection += 1
			elif re.match ( "^[0-9a-z-_]*$", str ( char ), re.IGNORECASE ):
				filter += str ( char )
				selection = -2
		ANSI.clear ()
		ANSI.cnorm ()
		return return_value

def main ():
	ANSI.smcup ()
	try:
		kube = Kube ()
		chosen_context = kube.prompt (
			"Pick Kubernetes Context",
			kube.context_names,
			kube.active_context_index
		)
		if chosen_context != None:
			kube.set_context ( chosen_context )
			chosen_namespace = kube.prompt (
				"Pick Kubernetes Namespace",
				kube.namespace_names,
				kube.active_namespace_index
			)
			if chosen_namespace != None:
				kube.set_namespace ( chosen_namespace )
	except ( KeyboardInterrupt, SystemExit ):
		ANSI.rmcup ()
		ANSI.cnorm ()
	except:
		ANSI.rmcup ()
		ANSI.cnorm ()
		print ( ANSI.error_color ("Error") + ": failed to communicate with k8s." )
		# raise
	else:
		ANSI.rmcup ()
		ANSI.cnorm ()

if __name__ == "__main__":
	main ()
