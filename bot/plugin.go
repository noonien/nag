package bot

import "errors"

type PluginInfo struct {
	Name        string
	Author      string
	Description string
	Version     string
}

type Plugin interface {
	Load(bot *Bot) (*PluginInfo, error)
	Unload() error
}

type botPlugin struct {
	*Bot

	plugin Plugin
	Info   PluginInfo

	handlerNames []string
	handlers     []Handler
}

func (bp *botPlugin) Handle(name string, handler Handler) {
	bp.handlerNames = append(bp.handlerNames, name)
	bp.handlers = append(bp.handlers, handler)
	bp.Bot.Handle(name, handler)
}

func (bp *botPlugin) Unload() error {
	for i := range bp.handlerNames {
		bp.Bot.RemoveHandler(bp.handlerNames[i], bp.handlers[i])
	}

	return bp.plugin.Unload()
}

func (b *Bot) LoadPlugin(plugin Plugin) error {
	info, err := plugin.Load(b)
	if err != nil {
		return err
	}

	bp := &botPlugin{
		Bot:    b,
		plugin: plugin,
		Info:   *info,
	}

	b.plugins[info.Name] = bp
	return nil
}

func (b *Bot) UnloadPlugin(name string) error {
	bp, ok := b.plugins[name]
	if !ok {
		return errors.New("plugin not loaded")
	}

	return bp.Unload()
}
