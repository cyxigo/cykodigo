package bot

import "github.com/bwmarrin/discordgo"

const (
	CmdHelp     = "help"
	CmdMeow     = "meow"
	CmdMeowat   = "meowat"
	CmdBark     = "bark"
	CmdBarkat   = "barkat"
	CmdRoulette = "roulette"
	CmdMe       = "me"
	CmdAssault  = "assault"
	CmdSexnkill = "sexnkill"
	CmdCat      = "cat"
	CmdCart     = "cart"
	CmdDoflip   = "doflip"
	CmdWork     = "work"
	CmdBalance  = "balance"
	CmdTransfer = "transfer"
	CmdSteal    = "steal"

	// msg means that this is message content command
	// same for handlers: handleMsg
	CmdMsgMeow  = "meow"
	CmdMsgCrazy = "crazy"
)

// slash commands, looooooooooong list of them
var cmds = []*discordgo.ApplicationCommand{
	{
		Name:        CmdHelp,
		Description: "This... will not help actually",
	},
	{
		Name:        CmdMeow,
		Description: "He will meow",
	},
	{
		Name:        CmdMeowat,
		Description: "Meow at someone!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to meow at",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdBark,
		Description: "He will... or won't bark. I don't know he's a cat",
	},
	{
		Name:        CmdBarkat,
		Description: "Bark at someone!", Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to bark at",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdRoulette,
		Description: "Why don't we play a little game?",
	},
	{
		Name:        CmdMe,
		Description: "Send picture of me! Nah, not you of course",
	},
	{
		Name:        CmdAssault,
		Description: "Try to assault someone... shh...",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to try to assault",
				Required:    true,
			},
		},
	},
	{
		// !!! its a joke command !!!
		Name:        CmdSexnkill,
		Description: "I don't know why you would do that.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to... you know",
				Required:    true,
			},
		},
	},
	{
		Name:        CmdCat,
		Description: "Cat!",
	},
	{
		Name:        CmdCart,
		Description: "Cart! (totally not copied from Cat Bot hehe)",
	},
	{
		Name:        CmdDoflip,
		Description: "He will do a flip",
	},
	{
		Name:        CmdWork,
		Description: "Work and get paid! Money!1!11!!",
	},
	{
		Name:        CmdBalance,
		Description: "Shows your or someone's money balance, so democratic!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person whose balance to show (optional)",
				Required:    false,
			},
		},
	},
	{
		Name:        CmdTransfer,
		Description: "Transfer your money to another person",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to whom to transfer",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of money to transfer",
				Required:    true,
			},
		},
	},
	{
		Name:        "steal",
		Description: "Steal money from someone! You can fail though, be careful.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "member",
				Description: "Person to steal money from",
				Required:    true,
			},
		},
	},
}
